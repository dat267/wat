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
	function getBarColor(val: number) {
		if (val > 85) return 'bg-rose-500/80';
		if (val > 65) return 'bg-amber-500/80';
		return 'bg-blue-500/80';
	}
</script>
<div class="space-y-8 select-none max-w-5xl">
	<div>
		<h1 class="text-2xl font-bold tracking-tight text-white mb-1">System Metrics</h1>
		<p class="text-slate-400 text-xs">Real-time system telemetry and environment resource utilization.</p>
	</div>
	{#if error}
		<div class="text-rose-400 text-xs font-semibold bg-rose-500/5 border border-rose-500/10 px-3 py-2 rounded-xl">{error}</div>
	{/if}
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6 transition-all duration-300 hover:border-white/[0.06] hover:bg-white/[0.02]">
			<div class="flex justify-between items-center mb-5">
				<span class="text-xs font-semibold text-slate-300 tracking-wider uppercase font-mono">CPU Load</span>
				<Cpu size={14} class="text-slate-400" />
			</div>
			<div class="flex items-baseline gap-1 mb-4">
				<span class="text-3xl font-light tracking-tight text-white">{stats ? stats.cpuPercent.toFixed(1) : '0.0'}</span>
				<span class="text-xs text-slate-400 font-semibold">%</span>
			</div>
			<div class="h-[3px] bg-white/[0.03] rounded-full overflow-hidden">
				<div
					class="h-full rounded-full transition-all duration-1000 {stats ? getBarColor(stats.cpuPercent) : 'bg-blue-500/80'}"
					style="width: {stats ? Math.min(stats.cpuPercent, 100) : 0}%"
				></div>
			</div>
		</div>
		<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6 transition-all duration-300 hover:border-white/[0.06] hover:bg-white/[0.02]">
			<div class="flex justify-between items-center mb-5">
				<span class="text-xs font-semibold text-slate-300 tracking-wider uppercase font-mono">Memory</span>
				<HardDrive size={14} class="text-slate-400" />
			</div>
			<div class="flex items-baseline gap-1 mb-4">
				<span class="text-3xl font-light tracking-tight text-white">{stats ? stats.memoryPercent.toFixed(1) : '0.0'}</span>
				<span class="text-xs text-slate-400 font-semibold">%</span>
			</div>
			<div class="h-[3px] bg-white/[0.03] rounded-full overflow-hidden mb-3">
				<div
					class="h-full rounded-full transition-all duration-1000 {stats ? getBarColor(stats.memoryPercent) : 'bg-blue-500/80'}"
					style="width: {stats ? Math.min(stats.memoryPercent, 100) : 0}%"
				></div>
			</div>
			<div class="text-xs text-slate-300 font-mono">
				{stats ? `${stats.memoryUsed.toFixed(1)} GB of ${stats.memoryTotal.toFixed(1)} GB` : '-'}
			</div>
		</div>
		<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6 transition-all duration-300 hover:border-white/[0.06] hover:bg-white/[0.02]">
			<div class="flex justify-between items-center mb-5">
				<span class="text-xs font-semibold text-slate-300 tracking-wider uppercase font-mono">Storage Capacity</span>
				<HardDrive size={14} class="text-slate-400" />
			</div>
			<div class="flex items-baseline gap-1 mb-4">
				<span class="text-3xl font-light tracking-tight text-white">{stats ? stats.diskPercent.toFixed(1) : '0.0'}</span>
				<span class="text-xs text-slate-400 font-semibold">%</span>
			</div>
			<div class="h-[3px] bg-white/[0.03] rounded-full overflow-hidden mb-3">
				<div
					class="h-full rounded-full transition-all duration-1000 {stats ? getBarColor(stats.diskPercent) : 'bg-blue-500/80'}"
					style="width: {stats ? Math.min(stats.diskPercent, 100) : 0}%"
				></div>
			</div>
			<div class="text-xs text-slate-300 font-mono">
				{stats ? `${stats.diskUsed.toFixed(1)} GB of ${stats.diskTotal.toFixed(1)} GB` : '-'}
			</div>
		</div>
	</div>
	<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-5 flex items-center gap-4 transition-all hover:border-white/[0.06]">
		<div class="p-2.5 rounded-xl bg-slate-900/50 text-slate-400">
			<Clock size={16} />
		</div>
		<div>
			<h3 class="text-xs font-semibold text-white">System Uptime</h3>
			<p class="text-slate-300 text-xs mt-0.5 font-mono">
				Continuous uptime: <strong class="text-white ml-1">{stats ? stats.uptime : 'Loading...'}</strong>
			</p>
		</div>
	</div>
</div>
