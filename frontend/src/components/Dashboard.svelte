<script lang="ts">
	import { Cpu, HardDrive, Clock } from '@lucide/svelte';
	import { GetSystemStats } from '../../wailsjs/go/main/App';
	import type { main } from '../../wailsjs/go/models';
	let stats = $state<main.SystemStats | null>(null);
	let error = $state<string | null>(null);
	$effect(() => {
		const fetchStats = async () => {
			try {
				const data = await GetSystemStats();
				stats = data;
				error = null;
			} catch (err) {
				error = 'Failed to fetch system metrics';
			}
		};
		fetchStats();
		const interval = setInterval(fetchStats, 2000);
		return () => clearInterval(interval);
	});
	function getProgressBarColor(val: number) {
		if (val > 85) return 'bg-rose-500 shadow-rose-500/40 text-rose-500';
		if (val > 65) return 'bg-amber-500 shadow-amber-500/40 text-amber-500';
		return 'bg-blue-500 shadow-blue-500/40 text-blue-500';
	}
</script>
<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-extrabold tracking-tight text-white mb-1.5">System Dashboard</h1>
		<p class="text-slate-400 text-sm">Real-time monitoring of local environment and hardware resource metrics.</p>
	</div>
	{#if error}
		<div class="text-rose-400 text-sm font-medium">{error}</div>
	{/if}
	<div class="grid grid-cols-1 md:grid-cols-3 gap-5">
		<div class="bg-slate-900/40 backdrop-blur-md border border-white/5 rounded-2xl p-6 transition-all duration-300 hover:-translate-y-1 hover:border-blue-500/20 hover:shadow-xl hover:shadow-blue-500/5">
			<div class="flex justify-between items-center mb-4">
				<span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">CPU UTILIZATION</span>
				<div class="p-2.5 rounded-xl bg-white/5 text-blue-500">
					<Cpu size={18} />
				</div>
			</div>
			<div class="flex items-baseline gap-1.5 mb-4">
				<span class="text-4xl font-bold tracking-tight text-white">{stats ? stats.cpuPercent.toFixed(1) : '0.0'}</span>
				<span class="text-sm text-slate-500 font-semibold">%</span>
			</div>
			<div class="h-1.5 bg-white/5 rounded-full overflow-hidden">
				<div
					class="h-full rounded-full transition-all duration-1000 shadow-[0_0_10px] {stats ? getProgressBarColor(stats.cpuPercent) : 'bg-blue-500 shadow-blue-500/40'}"
					style="width: {stats ? Math.min(stats.cpuPercent, 100) : 0}%"
				></div>
			</div>
		</div>
		<div class="bg-slate-900/40 backdrop-blur-md border border-white/5 rounded-2xl p-6 transition-all duration-300 hover:-translate-y-1 hover:border-blue-500/20 hover:shadow-xl hover:shadow-blue-500/5">
			<div class="flex justify-between items-center mb-4">
				<span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">SYSTEM MEMORY</span>
				<div class="p-2.5 rounded-xl bg-white/5 text-emerald-500">
					<HardDrive size={18} />
				</div>
			</div>
			<div class="flex items-baseline gap-1.5 mb-4">
				<span class="text-4xl font-bold tracking-tight text-white">{stats ? stats.memoryPercent.toFixed(1) : '0.0'}</span>
				<span class="text-sm text-slate-500 font-semibold">%</span>
			</div>
			<div class="h-1.5 bg-white/5 rounded-full overflow-hidden mb-3">
				<div
					class="h-full rounded-full transition-all duration-1000 shadow-[0_0_10px] {stats ? getProgressBarColor(stats.memoryPercent) : 'bg-emerald-500 shadow-emerald-500/40'}"
					style="width: {stats ? Math.min(stats.memoryPercent, 100) : 0}%"
				></div>
			</div>
			<div class="text-[11px] text-slate-400 font-medium">
				{stats ? `${stats.memoryUsed.toFixed(1)} GB used of ${stats.memoryTotal.toFixed(1)} GB` : '-'}
			</div>
		</div>
		<div class="bg-slate-900/40 backdrop-blur-md border border-white/5 rounded-2xl p-6 transition-all duration-300 hover:-translate-y-1 hover:border-blue-500/20 hover:shadow-xl hover:shadow-blue-500/5">
			<div class="flex justify-between items-center mb-4">
				<span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">DISK CAPACITY (ROOT)</span>
				<div class="p-2.5 rounded-xl bg-white/5 text-violet-500">
					<HardDrive size={18} />
				</div>
			</div>
			<div class="flex items-baseline gap-1.5 mb-4">
				<span class="text-4xl font-bold tracking-tight text-white">{stats ? stats.diskPercent.toFixed(1) : '0.0'}</span>
				<span class="text-sm text-slate-500 font-semibold">%</span>
			</div>
			<div class="h-1.5 bg-white/5 rounded-full overflow-hidden mb-3">
				<div
					class="h-full rounded-full transition-all duration-1000 shadow-[0_0_10px] {stats ? getProgressBarColor(stats.diskPercent) : 'bg-violet-500 shadow-violet-500/40'}"
					style="width: {stats ? Math.min(stats.diskPercent, 100) : 0}%"
				></div>
			</div>
			<div class="text-[11px] text-slate-400 font-medium">
				{stats ? `${stats.diskUsed.toFixed(1)} GB used of ${stats.diskTotal.toFixed(1)} GB` : '-'}
			</div>
		</div>
	</div>
	<div class="bg-slate-900/30 backdrop-blur-md border border-white/5 rounded-2xl p-5 flex items-center gap-4 transition-all hover:border-white/10">
		<div class="p-3 rounded-xl bg-amber-500/10 text-amber-500">
			<Clock size={22} />
		</div>
		<div>
			<h3 class="text-sm font-semibold text-white">Host System Uptime</h3>
			<p class="text-slate-400 text-xs mt-0.5">
				Operating continuous duration: <strong class="text-white ml-1 font-mono">{stats ? stats.uptime : 'Loading...'}</strong>
			</p>
		</div>
	</div>
</div>
