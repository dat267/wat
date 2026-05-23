<script lang="ts">
	import { LayoutDashboard, Braces, Network, Settings as SettingsIcon } from '@lucide/svelte';
	import Dashboard from './components/Dashboard.svelte';
	import DevTools from './components/DevTools.svelte';
	import PortScanner from './components/PortScanner.svelte';
	import Settings from './components/Settings.svelte';
	import { initTheme } from './theme.svelte';
	type Tab = 'dashboard' | 'tools' | 'ports' | 'settings';
	let activeTab = $state<Tab>('dashboard');
	initTheme();
</script>
<div class="flex h-screen w-screen bg-slate-950 text-slate-100 overflow-hidden font-sans select-none">
	<aside class="w-60 bg-slate-900/10 backdrop-blur-3xl border-r border-white/[0.04] flex flex-col p-6 z-10 justify-between">
		<div class="space-y-8">
			<div class="flex items-center gap-2.5 px-2">
				<div class="h-2 w-2 rounded-full bg-blue-500 animate-pulse"></div>
				<span class="text-xs font-bold tracking-widest text-slate-300 uppercase font-mono">WAT TOOLBOX</span>
			</div>
			<nav class="flex flex-col gap-1">
				<button
					class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs font-medium cursor-pointer transition-all border text-left w-full {activeTab === 'dashboard' ? 'text-white bg-white/[0.03] border-white/[0.05] font-semibold' : 'text-slate-300 border-transparent hover:text-white hover:bg-white/[0.01]'}"
					onclick={() => activeTab = 'dashboard'}
				>
					<LayoutDashboard size={14} />
					<span>Dashboard</span>
				</button>
				<button
					class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs font-medium cursor-pointer transition-all border text-left w-full {activeTab === 'tools' ? 'text-white bg-white/[0.03] border-white/[0.05] font-semibold' : 'text-slate-300 border-transparent hover:text-white hover:bg-white/[0.01]'}"
					onclick={() => activeTab = 'tools'}
				>
					<Braces size={14} />
					<span>Disk Ranker</span>
				</button>
				<button
					class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs font-medium cursor-pointer transition-all border text-left w-full {activeTab === 'ports' ? 'text-white bg-white/[0.03] border-white/[0.05] font-semibold' : 'text-slate-300 border-transparent hover:text-white hover:bg-white/[0.01]'}"
					onclick={() => activeTab = 'ports'}
				>
					<Network size={14} />
					<span>Port Scanner</span>
				</button>
				<button
					class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs font-medium cursor-pointer transition-all border text-left w-full {activeTab === 'settings' ? 'text-white bg-white/[0.03] border-white/[0.05] font-semibold' : 'text-slate-300 border-transparent hover:text-white hover:bg-white/[0.01]'}"
					onclick={() => activeTab = 'settings'}
				>
					<SettingsIcon size={14} />
					<span>Settings</span>
				</button>
			</nav>
		</div>
		<div class="text-[9px] text-slate-500 pt-4 border-t border-white/[0.03] flex justify-between tracking-widest font-mono">
			<span>v1.0.0</span>
			<span class="text-blue-500/80">CONNECTED</span>
		</div>
	</aside>
	<main class="flex-1 h-full overflow-y-auto p-8 bg-gradient-to-b from-slate-950 via-slate-950 to-slate-950/95">
		{#if activeTab === 'dashboard'}
			<Dashboard />
		{:else if activeTab === 'tools'}
			<DevTools />
		{:else if activeTab === 'ports'}
			<PortScanner />
		{:else if activeTab === 'settings'}
			<Settings />
		{/if}
	</main>
</div>
