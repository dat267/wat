<script lang="ts">
	import { LayoutDashboard, Braces, Settings as SettingsIcon } from '@lucide/svelte';
	import Dashboard from './components/Dashboard.svelte';
	import DevTools from './components/DevTools.svelte';
	import Settings from './components/Settings.svelte';
	import { initTheme } from './theme.svelte';
	type Tab = 'dashboard' | 'tools' | 'settings';
	let activeTab = $state<Tab>('dashboard');
	initTheme();
</script>
<div class="flex h-screen w-screen bg-slate-950 text-slate-100 overflow-hidden font-sans">
	<aside class="w-64 bg-slate-900/60 backdrop-blur-xl border-r border-white/5 flex flex-col p-6 z-10 justify-between">
		<div class="space-y-6">
			<div class="text-[10px] font-bold text-slate-500 tracking-wider pl-2 uppercase">Main Navigation</div>
			<nav class="flex flex-col gap-1.5">
				<button
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-400 text-sm font-medium cursor-pointer transition-all border border-transparent hover:text-white hover:bg-white/5 text-left w-full {activeTab === 'dashboard' ? 'text-white bg-blue-500/10 border-blue-500/25 shadow-sm' : ''}"
					onclick={() => activeTab = 'dashboard'}
				>
					<LayoutDashboard size={18} />
					<span>Metrics Dashboard</span>
				</button>
				<button
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-400 text-sm font-medium cursor-pointer transition-all border border-transparent hover:text-white hover:bg-white/5 text-left w-full {activeTab === 'tools' ? 'text-white bg-blue-500/10 border-blue-500/25 shadow-sm' : ''}"
					onclick={() => activeTab = 'tools'}
				>
					<Braces size={18} />
					<span>Tools</span>
				</button>
				<button
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-400 text-sm font-medium cursor-pointer transition-all border border-transparent hover:text-white hover:bg-white/5 text-left w-full {activeTab === 'settings' ? 'text-white bg-blue-500/10 border-blue-500/25 shadow-sm' : ''}"
					onclick={() => activeTab = 'settings'}
				>
					<SettingsIcon size={18} />
					<span>Settings</span>
				</button>
			</nav>
		</div>
		<div class="text-[10px] text-slate-500 pt-4 border-t border-white/5 flex justify-between tracking-wider font-semibold">
			<span>VERSION 1.0.0</span>
			<span class="text-emerald-500">SECURE</span>
		</div>
	</aside>
	<main class="flex-1 h-full overflow-y-auto p-8 bg-gradient-to-b from-slate-900/10 to-slate-950">
		{#if activeTab === 'dashboard'}
			<Dashboard />
		{:else if activeTab === 'tools'}
			<DevTools />
		{:else if activeTab === 'settings'}
			<Settings />
		{/if}
	</main>
</div>
