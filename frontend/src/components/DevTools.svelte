<script lang="ts">
	import { HardDrive, FolderOpen, File, Search, RefreshCw } from '@lucide/svelte';
	import { RankDirectory } from '../../wailsjs/go/main/App';
	import { toastState } from './toast.svelte';
	import type { main } from '../../wailsjs/go/models';
	let path = $state('.');
	let entries = $state<main.DiskEntry[]>([]);
	let filterQuery = $state('');
	let total = $state(0);
	let loading = $state(false);
	let error = $state<string | null>(null);
	let scanned = $state(false);
	const filteredEntries = $derived(
		entries.filter(e => e.name.toLowerCase().includes(filterQuery.toLowerCase()))
	);
	async function scan(targetPath = path) {
		if (!targetPath.trim()) return;
		path = targetPath;
		loading = true;
		error = null;
		scanned = false;
		try {
			entries = await RankDirectory(targetPath.trim());
			total = entries.reduce((s, e) => s + e.size, 0);
			scanned = true;
			toastState.add('success', `Directory ranked: ${entries.length} items analyzed.`);
		} catch (e: any) {
			error = e?.toString() ?? 'Failed to scan directory';
			toastState.add('error', 'Failed to rank directory.');
		} finally {
			loading = false;
		}
	}
	function humanSize(b: number): string {
		if (b >= 1073741824) return (b / 1073741824).toFixed(1) + ' GB';
		if (b >= 1048576) return (b / 1048576).toFixed(1) + ' MB';
		if (b >= 1024) return (b / 1024).toFixed(1) + ' KB';
		return b + ' B';
	}
	function barWidth(size: number): string {
		if (total === 0) return '0%';
		return Math.min((size / total) * 100, 100).toFixed(2) + '%';
	}
	function pct(size: number): string {
		if (total === 0) return '0.0';
		return ((size / total) * 100).toFixed(1);
	}
	function barColor(size: number): string {
		const p = total > 0 ? (size / total) * 100 : 0;
		if (p > 40) return 'bg-rose-500/60';
		if (p > 20) return 'bg-amber-500/60';
		return 'bg-blue-500/60';
	}
</script>
<div class="space-y-8 select-none max-w-5xl">
	<div>
		<h1 class="text-2xl font-bold tracking-tight text-white mb-1">Disk Capacity Ranker</h1>
		<p class="text-slate-400 text-xs">Analyze and rank folders and files by storage usage metrics.</p>
	</div>
	<div class="space-y-3">
		<div class="flex gap-2">
			<div class="relative flex-1">
				<div class="absolute inset-y-0 left-3 flex items-center pointer-events-none text-slate-500">
					<FolderOpen size={14} />
				</div>
				<input
					id="disk-path-input"
					type="text"
					bind:value={path}
					placeholder="e.g. /home/user"
					onkeydown={(e) => e.key === 'Enter' && scan()}
					class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl pl-9 pr-4 py-2.5 text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-white/[0.08] transition-all font-mono"
				/>
			</div>
			<button
				id="disk-scan-btn"
				onclick={() => scan()}
				disabled={loading || !path.trim()}
				class="px-5 py-2.5 rounded-xl bg-white text-slate-950 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-200 text-xs font-semibold transition-all active:scale-95 cursor-pointer font-mono"
			>
				{loading ? 'ANALYZING…' : 'ANALYZE'}
			</button>
		</div>
		<div class="flex items-center gap-2 text-[10px] font-mono text-slate-400">
			<span>Quick paths:</span>
			<button onclick={() => scan('.')} class="hover:text-white transition-colors underline cursor-pointer">current dir (.)</button>
			<span>•</span>
			<button onclick={() => scan('/')} class="hover:text-white transition-colors underline cursor-pointer">root (/)</button>
		</div>
	</div>
	{#if error}
		<div class="text-rose-400 text-xs font-semibold bg-rose-500/5 border border-rose-500/10 px-4 py-3 rounded-xl">{error}</div>
	{/if}
	{#if scanned}
		{#if entries.length === 0}
			<div class="text-slate-400 text-xs text-center py-12">Directory is empty.</div>
		{:else}
			<div class="space-y-4">
				<div class="flex items-center justify-between gap-4">
					<div class="relative max-w-xs w-full">
						<div class="absolute inset-y-0 left-3 flex items-center pointer-events-none text-slate-500">
							<Search size={12} />
						</div>
						<input
							type="text"
							bind:value={filterQuery}
							placeholder="Filter list by name..."
							class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl pl-8 pr-4 py-2 text-[11px] text-white focus:outline-none focus:border-white/[0.08] font-mono"
						/>
					</div>
					<div class="text-[11px] text-slate-400 font-mono">
						Showing {filteredEntries.length} of {entries.length} items
					</div>
				</div>
				<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl overflow-hidden">
					<div class="flex items-center justify-between px-5 py-3 border-b border-white/[0.03] bg-white/[0.005]">
						<div class="flex items-center gap-2 text-slate-300 text-xs font-semibold uppercase tracking-wider font-mono">
							<HardDrive size={12} />
							<span>{path}</span>
						</div>
						<span class="text-slate-300 text-xs font-mono">{humanSize(total)} total</span>
					</div>
					<div class="divide-y divide-white/[0.02] max-h-[450px] overflow-y-auto">
						{#each filteredEntries as entry (entry.path)}
							<div class="group flex items-center gap-4 px-5 py-3 hover:bg-white/[0.005] transition-colors">
								<div class="shrink-0 transition-colors">
									{#if entry.isDir}
										<span class="text-blue-400"><FolderOpen size={13} /></span>
									{:else}
										<span class="text-slate-500"><File size={13} /></span>
									{/if}
								</div>
								<div class="flex-1 min-w-0">
									<div class="flex items-center justify-between mb-1.5">
										<span class="text-xs text-white font-mono truncate pr-4">
											{entry.name}{entry.isDir ? '/' : ''}
										</span>
										<span class="text-xs text-slate-300 font-mono shrink-0">{humanSize(entry.size)}</span>
									</div>
									<div class="flex items-center gap-2.5">
										<div class="flex-1 h-[2px] bg-white/[0.03] rounded-full overflow-hidden">
											<div
												class="h-full rounded-full transition-all duration-500 {barColor(entry.size)}"
												style="width: {barWidth(entry.size)}"
											></div>
										</div>
										<span class="text-xs text-slate-400 font-mono w-10 text-right">{pct(entry.size)}%</span>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>
