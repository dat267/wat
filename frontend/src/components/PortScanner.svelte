<script lang="ts">
	import { Network, Play, ShieldAlert, Copy, RefreshCw, ChevronDown, ChevronUp } from '@lucide/svelte';
	import { ScanPorts } from '../../wailsjs/go/main/App';
	let host = $state('localhost');
	let scanType = $state<'common' | 'all' | 'custom'>('common');
	let customPortsStr = $state('1-1024');
	let timeoutMs = $state(500);
	let concurrency = $state(1024);
	let showAdvanced = $state(false);
	let scanning = $state(false);
	let openPorts = $state<number[]>([]);
	let error = $state<string | null>(null);
	let scanned = $state(false);
	let copied = $state(false);
	let dropdownOpen = $state(false);
	const serviceMap: Record<number, string> = {
		21: 'FTP',
		22: 'SSH',
		23: 'Telnet',
		25: 'SMTP',
		53: 'DNS',
		80: 'HTTP',
		110: 'POP3',
		143: 'IMAP',
		443: 'HTTPS',
		3306: 'MySQL',
		5432: 'PostgreSQL',
		6379: 'Redis',
		8080: 'HTTP-Alt',
		8443: 'HTTPS-Alt',
		27017: 'MongoDB'
	};
	function getServiceName(port: number): string {
		return serviceMap[port] ?? 'Unknown Service';
	}
	async function startScan() {
		if (!host.trim()) return;
		scanning = true;
		error = null;
		scanned = false;
		openPorts = [];
		try {
			let portsToScan: number[] = [];
			if (scanType === 'all') {
				portsToScan = Array.from({ length: 65535 }, (_, i) => i + 1);
			} else if (scanType === 'common') {
				portsToScan = [21, 22, 23, 25, 53, 80, 110, 143, 443, 3306, 5432, 6379, 8080, 8443, 27017];
			} else {
				portsToScan = parsePortsString(customPortsStr);
				if (portsToScan.length === 0) {
					throw new Error('No valid ports specified in custom ports list.');
				}
			}
			openPorts = await ScanPorts(host.trim(), portsToScan, timeoutMs, concurrency);
			scanned = true;
		} catch (e: any) {
			error = e?.toString() ?? 'Failed to perform port scan';
		} finally {
			scanning = false;
		}
	}
	function parsePortsString(str: string): number[] {
		const ports: number[] = [];
		const seen = new Set<number>();
		const parts = str.split(',');
		for (let part of parts) {
			part = part.trim();
			if (!part) continue;
			if (part.includes('-')) {
				const [startStr, endStr] = part.split('-');
				const start = parseInt(startStr, 10);
				const end = parseInt(endStr, 10);
				if (!isNaN(start) && !isNaN(end) && start <= end && start >= 1 && end <= 65535) {
					for (let i = start; i <= end; i++) {
						if (!seen.has(i)) {
							seen.add(i);
							ports.push(i);
						}
					}
				}
			} else {
				const p = parseInt(part, 10);
				if (!isNaN(p) && p >= 1 && p <= 65535 && !seen.has(p)) {
					seen.add(p);
					ports.push(p);
				}
			}
		}
		return ports.sort((a, b) => a - b);
	}
	function copyPorts() {
		if (openPorts.length === 0) return;
		navigator.clipboard.writeText(openPorts.join(', '));
		copied = true;
		setTimeout(() => copied = false, 2000);
	}
</script>
<div class="space-y-8 select-none max-w-5xl">
	<div>
		<h1 class="text-2xl font-bold tracking-tight text-white mb-1">Port Scanner</h1>
		<p class="text-slate-400 text-xs">Verify target network host availability and audit active port mappings.</p>
	</div>
	<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6 space-y-5">
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div class="space-y-1.5">
				<label for="ports-host-input" class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Target Host</label>
				<input
					id="ports-host-input"
					type="text"
					bind:value={host}
					disabled={scanning}
					placeholder="localhost"
					class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-4 py-2.5 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono transition-all"
				/>
			</div>
			<div class="space-y-1.5 relative">
				<span class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Range</span>
				{#if dropdownOpen}
					<div class="fixed inset-0 z-40" onclick={() => dropdownOpen = false}></div>
				{/if}
				<div class="relative z-50">
					<button
						type="button"
						onclick={() => !scanning && (dropdownOpen = !dropdownOpen)}
						disabled={scanning}
						class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-4 py-2.5 text-xs text-slate-300 focus:outline-none focus:border-white/[0.08] transition-all flex items-center justify-between cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
					>
						<span>
							{#if scanType === 'common'}
								Common Ports
							{:else}
								{scanType === 'all' ? 'All Ports (1 - 65535)' : 'Custom Configuration'}
							{/if}
						</span>
						{#if dropdownOpen}
							<ChevronUp size={14} class="text-slate-500" />
						{:else}
							<ChevronDown size={14} class="text-slate-500" />
						{/if}
					</button>
					{#if dropdownOpen}
						<div class="absolute mt-1.5 w-full bg-[#121726]/95 border border-white/[0.05] rounded-xl shadow-xl backdrop-blur-2xl p-1 flex flex-col gap-0.5">
							<button
								type="button"
								onclick={() => { scanType = 'common'; dropdownOpen = false; }}
								class="flex flex-col text-left px-3.5 py-2 rounded-lg text-xs transition-colors cursor-pointer hover:bg-white/[0.02] {scanType === 'common' ? 'bg-white/[0.03] text-white font-semibold' : 'text-slate-400'}"
							>
								<span>Common Ports</span>
								<span class="text-[9px] text-slate-500 font-normal mt-0.5 font-mono">15 standard operational ports</span>
							</button>
							<button
								type="button"
								onclick={() => { scanType = 'all'; dropdownOpen = false; }}
								class="flex flex-col text-left px-3.5 py-2 rounded-lg text-xs transition-colors cursor-pointer hover:bg-white/[0.02] {scanType === 'all' ? 'bg-white/[0.03] text-white font-semibold' : 'text-slate-400'}"
							>
								<span>All Ports</span>
								<span class="text-[9px] text-slate-500 font-normal mt-0.5 font-mono">Scan full range (1 - 65535)</span>
							</button>
							<button
								type="button"
								onclick={() => { scanType = 'custom'; dropdownOpen = false; }}
								class="flex flex-col text-left px-3.5 py-2 rounded-lg text-xs transition-colors cursor-pointer hover:bg-white/[0.02] {scanType === 'custom' ? 'bg-white/[0.03] text-white font-semibold' : 'text-slate-400'}"
							>
								<span>Custom Configuration</span>
								<span class="text-[9px] text-slate-500 font-normal mt-0.5 font-mono">Define custom ports or ranges</span>
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
		{#if scanType === 'custom'}
			<div class="space-y-1.5">
				<label for="ports-custom-input" class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Ports list</label>
				<input
					id="ports-custom-input"
					type="text"
					bind:value={customPortsStr}
					disabled={scanning}
					placeholder="e.g. 22,80,443,8080-8090"
					class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-4 py-2.5 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono transition-all"
				/>
			</div>
		{/if}
		<div class="border-t border-white/[0.03] pt-4">
			<button
				onclick={() => showAdvanced = !showAdvanced}
				class="flex items-center gap-1 text-xs font-semibold tracking-wider uppercase text-slate-400 hover:text-slate-200 transition-colors cursor-pointer"
			>
				{#if showAdvanced}
					<ChevronUp size={12} /> Hide Advanced Options
				{:else}
					<ChevronDown size={12} /> Show Advanced Options
				{/if}
			</button>
			{#if showAdvanced}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
					<div class="space-y-1.5">
						<label for="ports-timeout-input" class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Timeout (ms)</label>
						<input
							id="ports-timeout-input"
							type="number"
							bind:value={timeoutMs}
							disabled={scanning}
							min="50"
							max="5000"
							class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-4 py-2 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono"
						/>
					</div>
					<div class="space-y-1.5">
						<label for="ports-concurrency-input" class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Concurrency Workers</label>
						<input
							id="ports-concurrency-input"
							type="number"
							bind:value={concurrency}
							disabled={scanning}
							min="1"
							max="4096"
							class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-4 py-2 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono"
						/>
					</div>
				</div>
			{/if}
		</div>
		<div class="flex justify-end pt-2">
			<button
				id="ports-scan-btn"
				onclick={startScan}
				disabled={scanning || !host.trim()}
				class="px-6 py-2.5 rounded-xl bg-white text-slate-950 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-200 text-xs font-semibold flex items-center gap-2 transition-all active:scale-95 cursor-pointer font-mono"
			>
				{#if scanning}
					<RefreshCw size={12} class="animate-spin" /> SCANNING…
				{:else}
					<Play size={12} /> SCAN HOST
				{/if}
			</button>
		</div>
	</div>
	{#if error}
		<div class="text-rose-400 text-xs font-semibold bg-rose-500/5 border border-rose-500/10 px-4 py-3 rounded-xl">{error}</div>
	{/if}
	{#if scanned}
		<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl overflow-hidden">
			<div class="flex items-center justify-between px-5 py-3 border-b border-white/[0.03] bg-white/[0.005]">
				<span class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Open ports for {host}</span>
				<div class="flex items-center gap-3">
					<span class="text-slate-300 text-xs font-mono">{openPorts.length} open</span>
					{#if openPorts.length > 0}
						<button
							onclick={copyPorts}
							class="text-xs text-slate-300 hover:text-white flex items-center gap-1 cursor-pointer transition-colors font-mono"
						>
							<Copy size={11} /> {copied ? 'COPIED!' : 'COPY'}
						</button>
					{/if}
				</div>
			</div>
			{#if openPorts.length === 0}
				<div class="text-slate-400 text-xs text-center py-12">No open ports found.</div>
			{:else}
				<div class="divide-y divide-white/[0.02]">
					{#each openPorts as p}
						<div class="flex items-center justify-between px-5 py-3 hover:bg-white/[0.005] transition-colors">
							<div class="flex items-baseline gap-4">
								<span class="text-xs font-bold text-white font-mono w-16">{p}</span>
								<span class="text-xs text-slate-400 font-mono">{getServiceName(p)}</span>
							</div>
							<span class="px-1.5 py-0.5 rounded text-[8px] font-bold tracking-widest uppercase bg-blue-500/10 border border-blue-500/20 text-blue-400 font-mono">OPEN</span>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>
