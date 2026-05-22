<script lang="ts">
	import { Moon, Sun, Monitor, Info, ShieldCheck } from '@lucide/svelte';
	import { themeState, type Theme } from '../theme.svelte';
	const themes = [
		{ key: 'dark' as Theme, label: 'Dark Theme', icon: Moon, desc: 'Premium deep space obsidian look' },
		{ key: 'light' as Theme, label: 'Light Theme', icon: Sun, desc: 'Clean, crisp and highly readable' },
		{ key: 'system' as Theme, label: 'System Default', icon: Monitor, desc: 'Follows your operating system setting' }
	];
</script>
<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-extrabold tracking-tight text-white mb-1.5">Toolbox Settings</h1>
		<p class="text-slate-400 text-sm">Customize client workspace configuration, visual appearance, and internal features.</p>
	</div>
	<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
		<div class="bg-slate-900/40 border border-white/5 rounded-2xl p-6">
			<h3 class="text-sm font-semibold text-white flex items-center gap-2 mb-4">
				<Monitor size={16} class="text-blue-500" /> App Appearance Theme
			</h3>
			<p class="text-slate-400 text-xs leading-relaxed mb-4">
				Select your preferred display theme. This updates the global styling variables instantly:
			</p>
			<div class="flex flex-col gap-2.5">
				{#each themes as t}
					<button
						onclick={() => themeState.current = t.key}
						class="flex items-center gap-4 px-4 py-3.5 rounded-xl border cursor-pointer text-left transition-all {themeState.current === t.key ? 'border-blue-500 bg-blue-500/5 shadow-md shadow-blue-500/5 text-white' : 'border-white/5 bg-transparent text-slate-300 hover:bg-white/5'}"
					>
						<div class="p-2 rounded-lg {themeState.current === t.key ? 'bg-blue-500/20 text-blue-400' : 'bg-white/5 text-slate-400'}">
							<t.icon size={16} />
						</div>
						<div class="flex-1">
							<div class="text-xs font-semibold">{t.label}</div>
							<div class="text-[10px] text-slate-400 mt-0.5">{t.desc}</div>
						</div>
					</button>
				{/each}
			</div>
		</div>
		<div class="bg-slate-900/40 border border-white/5 rounded-2xl p-6">
			<h3 class="text-sm font-semibold text-white flex items-center gap-2 mb-4">
				<Info size={16} class="text-blue-500" /> Runtime Environment
			</h3>
			<div class="flex flex-col gap-3 text-xs">
				<div class="flex justify-between border-b border-white/5 pb-2.5">
					<span class="text-slate-400">Wails Platform</span>
					<strong class="text-slate-200 font-mono">v2.12.0</strong>
				</div>
				<div class="flex justify-between border-b border-white/5 pb-2.5">
					<span class="text-slate-400">Frontend Engine</span>
					<strong class="text-slate-200">Svelte 5 + Vite (TS)</strong>
				</div>
				<div class="flex justify-between border-b border-white/5 pb-2.5">
					<span class="text-slate-400">Backend Language</span>
					<strong class="text-slate-200 font-mono">Go (1.26.3)</strong>
				</div>
				<div class="flex justify-between pb-2.5">
					<span class="text-slate-400">Sandboxed Runtime</span>
					<strong class="text-emerald-400 flex items-center gap-1.5">
						<ShieldCheck size={14} /> Native OS Bindings
					</strong>
				</div>
			</div>
		</div>
	</div>
</div>
