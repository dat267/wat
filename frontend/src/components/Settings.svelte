<script lang="ts">
	import { Moon, Sun, Monitor, Info, ShieldCheck, Key, Eye, EyeOff, Save, CheckCircle, Plus, Trash2 } from '@lucide/svelte';
	import { themeState, type Theme } from '../theme.svelte';
	import { GetConfig, SetConfigValue, ResetConfig, DeleteConfigValue } from '../../wailsjs/go/main/App';
	const themes = [
		{ key: 'dark' as Theme, label: 'Dark', icon: Moon, desc: 'Obsidian space look' },
		{ key: 'light' as Theme, label: 'Light', icon: Sun, desc: 'Clean, highly readable' },
		{ key: 'system' as Theme, label: 'System', icon: Monitor, desc: 'Follows OS appearance' }
	];
	let openaiKey = $state('');
	let anthropicKey = $state('');
	let customKey = $state('');
	let customVal = $state('');
	let showOpenai = $state(false);
	let showAnthropic = $state(false);
	let loading = $state(true);
	let saving = $state(false);
	let savedMessage = $state<string | null>(null);
	let customKeysList = $state<{ key: string; value: string }[]>([]);
	let confirmDelete = $state(false);
	async function performResetConfig() {
		try {
			await ResetConfig();
			confirmDelete = false;
			openaiKey = '';
			anthropicKey = '';
			customKeysList = [];
			customKey = '';
			customVal = '';
			themeState.current = 'system';
			savedMessage = 'Configuration reset successfully!';
			setTimeout(() => savedMessage = null, 3000);
		} catch (e) {
			console.error('Failed to reset config file', e);
		}
	}
	$effect(() => {
		loadSecureConfig();
	});
	async function loadSecureConfig() {
		try {
			const cfg = await GetConfig();
			openaiKey = cfg.apiKeys['openai'] || '';
			anthropicKey = cfg.apiKeys['anthropic'] || '';
			customKeysList = Object.entries(cfg.apiKeys)
				.filter(([k]) => k !== 'openai' && k !== 'anthropic')
				.map(([k, v]) => ({ key: k, value: v }));
			if (cfg.settings['theme']) {
				themeState.current = cfg.settings['theme'] as Theme;
			}
		} catch (e) {
			console.error('Failed to load secure config', e);
		} finally {
			loading = false;
		}
	}
	async function saveKeys() {
		saving = true;
		savedMessage = null;
		try {
			await SetConfigValue('apiKeys', 'openai', openaiKey);
			await SetConfigValue('apiKeys', 'anthropic', anthropicKey);
			await SetConfigValue('settings', 'theme', themeState.current);
			savedMessage = 'Configuration saved securely!';
			setTimeout(() => savedMessage = null, 3000);
		} catch (e: any) {
			console.error('Failed to save config', e);
		} finally {
			saving = false;
		}
	}
	async function addCustomKey() {
		if (!customKey.trim() || !customVal.trim()) return;
		try {
			await SetConfigValue('apiKeys', customKey.trim(), customVal.trim());
			customKey = '';
			customVal = '';
			await loadSecureConfig();
		} catch (e) {
			console.error(e);
		}
	}
	async function removeCustomKey(k: string) {
		try {
			await DeleteConfigValue('apiKeys', k);
			await loadSecureConfig();
		} catch (e) {
			console.error(e);
		}
	}
	async function changeTheme(key: Theme) {
		themeState.current = key;
		try {
			await SetConfigValue('settings', 'theme', key);
		} catch (e) {
			console.error(e);
		}
	}
</script>
<div class="space-y-8 select-none max-w-5xl">
	<div>
		<h1 class="text-2xl font-bold tracking-tight text-white mb-1">Configuration Settings</h1>
		<p class="text-slate-400 text-xs">Manage local environment variables, visual preferences, and API credentials.</p>
	</div>
	{#if loading}
		<div class="text-slate-400 text-xs text-center py-12">Loading secure settings…</div>
	{:else}
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<div class="lg:col-span-2 space-y-6">
				<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6 space-y-5">
					<h3 class="text-xs font-semibold text-white flex items-center gap-2 pb-3 border-b border-white/[0.03]">
						<Key size={13} class="text-slate-400" />
						<span>Secure Credentials</span>
					</h3>
					<p class="text-slate-300 text-xs leading-relaxed">
						API keys are encrypted locally using owner-only permissions (<code class="font-mono bg-slate-950 px-1 rounded text-slate-300">0600</code>).
					</p>
					<div class="space-y-4">
						<div class="space-y-1.5">
							<div class="flex items-center justify-between">
								<label for="settings-openai-input" class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">OpenAI API Key</label>
								<span class="text-[8px] font-bold text-blue-400 uppercase tracking-widest px-1.5 py-0.5 rounded bg-blue-500/10 border border-blue-500/20 font-mono">Secret</span>
							</div>
							<div class="relative">
								<input
									id="settings-openai-input"
									type={showOpenai ? 'text' : 'password'}
									bind:value={openaiKey}
									placeholder="sk-..."
									class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl pl-4 pr-10 py-2.5 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono transition-all"
								/>
								<button
									onclick={() => showOpenai = !showOpenai}
									class="absolute right-3 top-3 text-slate-500 hover:text-white transition-colors cursor-pointer"
								>
									{#if showOpenai}
										<EyeOff size={13} />
									{:else}
										<Eye size={13} />
									{/if}
								</button>
							</div>
						</div>
						<div class="space-y-1.5">
							<div class="flex items-center justify-between">
								<label for="settings-anthropic-input" class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Anthropic API Key</label>
								<span class="text-[8px] font-bold text-blue-400 uppercase tracking-widest px-1.5 py-0.5 rounded bg-blue-500/10 border border-blue-500/20 font-mono">Secret</span>
							</div>
							<div class="relative">
								<input
									id="settings-anthropic-input"
									type={showAnthropic ? 'text' : 'password'}
									bind:value={anthropicKey}
									placeholder="sk-ant-..."
									class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl pl-4 pr-10 py-2.5 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono transition-all"
								/>
								<button
									onclick={() => showAnthropic = !showAnthropic}
									class="absolute right-3 top-3 text-slate-500 hover:text-white transition-colors cursor-pointer"
								>
									{#if showAnthropic}
										<EyeOff size={13} />
									{:else}
										<Eye size={13} />
									{/if}
								</button>
							</div>
						</div>
					</div>
					<div class="border-t border-white/[0.03] pt-4 space-y-3">
						<div class="flex items-center justify-between">
							<span class="text-xs font-semibold text-slate-300 uppercase tracking-wider font-mono">Custom Keys</span>
							<span class="text-[8px] font-bold text-blue-400 uppercase tracking-widest px-1.5 py-0.5 rounded bg-blue-500/10 border border-blue-500/20 font-mono">Secret</span>
						</div>
						<div class="flex gap-2">
							<input
								type="text"
								bind:value={customKey}
								placeholder="e.g. github_token"
								class="flex-1 bg-white/[0.01] border border-white/[0.04] rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono"
							/>
							<input
								type="password"
								bind:value={customVal}
								placeholder="Value"
								class="flex-1 bg-white/[0.01] border border-white/[0.04] rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-white/[0.08] font-mono"
							/>
							<button
								onclick={addCustomKey}
								disabled={!customKey.trim() || !customVal.trim()}
								class="px-3 rounded-xl bg-white text-slate-950 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-200 text-xs font-semibold flex items-center gap-1 transition-all cursor-pointer font-mono"
							>
								<Plus size={12} /> ADD
							</button>
						</div>
						{#if customKeysList.length > 0}
							<div class="bg-white/[0.005] border border-white/[0.03] rounded-xl divide-y divide-white/[0.02] max-h-40 overflow-y-auto font-mono">
								{#each customKeysList as item}
									<div class="flex items-center justify-between px-3 py-2.5">
										<span class="text-xs text-slate-300">{item.key}</span>
										<div class="flex items-center gap-3">
											<span class="text-xs text-slate-400">[encrypted]</span>
											<button
												onclick={() => removeCustomKey(item.key)}
												class="text-rose-500 hover:text-rose-400 p-1 cursor-pointer transition-colors"
											>
												<Trash2 size={12} />
											</button>
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
					<div class="flex items-center justify-between border-t border-white/[0.03] pt-4">
						{#if savedMessage}
							<div class="text-emerald-400 text-xs font-medium flex items-center gap-1.5">
								<CheckCircle size={12} /> {savedMessage}
							</div>
						{:else}
							<div></div>
						{/if}
						<button
							onclick={saveKeys}
							disabled={saving}
							class="px-5 py-2.5 rounded-xl bg-white text-slate-950 hover:bg-slate-200 disabled:opacity-40 text-xs font-semibold flex items-center gap-2 transition-all active:scale-95 cursor-pointer font-mono"
						>
							<Save size={12} /> SAVE CONFIG
						</button>
					</div>
				</div>
				<div class="bg-rose-950/5 border border-rose-500/10 rounded-2xl p-6 space-y-4">
					<h3 class="text-xs font-semibold text-rose-400 flex items-center gap-2.5 pb-3 border-b border-rose-500/10">
						<Trash2 size={13} class="text-rose-500" />
						<span>Danger Zone</span>
					</h3>
					<p class="text-slate-300 text-xs leading-relaxed">
						Reset all secure configurations and visual preferences to their standard defaults.
					</p>
					<div class="flex justify-end">
						{#if confirmDelete}
							<div class="flex items-center gap-2 font-mono">
								<span class="text-[10px] text-rose-400 font-medium">Are you sure?</span>
								<button
									onclick={performResetConfig}
									class="px-4 py-2 rounded-xl bg-rose-600 hover:bg-rose-500 text-white text-[10px] font-semibold cursor-pointer transition-all active:scale-95"
								>
									YES, RESET
								</button>
								<button
									onclick={() => confirmDelete = false}
									class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold cursor-pointer transition-all"
								>
									CANCEL
								</button>
							</div>
						{:else}
							<button
								onclick={() => confirmDelete = true}
								class="px-5 py-2.5 rounded-xl bg-rose-950/20 hover:bg-rose-900/20 border border-rose-500/20 text-rose-400 hover:text-rose-300 text-xs font-semibold cursor-pointer transition-all active:scale-95 font-mono"
							>
								RESET SECURE CONFIG TO DEFAULTS
							</button>
						{/if}
					</div>
				</div>
			</div>
			<div class="space-y-6">
				<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6">
					<h3 class="text-xs font-semibold text-white flex items-center justify-between mb-4">
						<span class="flex items-center gap-2">
							<Monitor size={13} class="text-slate-400" /> Display Theme
						</span>
						<span class="text-[8px] font-bold text-rose-400 uppercase tracking-widest px-1.5 py-0.5 rounded bg-rose-500/10 border border-rose-500/20 font-mono">Mandatory</span>
					</h3>
					<div class="flex flex-col gap-2">
						{#each themes as t}
							<button
								onclick={() => changeTheme(t.key)}
								class="flex items-center gap-4 px-4 py-3 rounded-xl border cursor-pointer text-left transition-all {themeState.current === t.key ? 'border-white/[0.08] bg-white/[0.03] text-white font-semibold' : 'border-transparent bg-transparent text-slate-400 hover:bg-white/[0.01]'}"
							>
								<div class="p-2 rounded-lg {themeState.current === t.key ? 'bg-white/[0.05] text-white' : 'bg-transparent text-slate-500'}">
									<t.icon size={13} />
								</div>
								<div class="flex-1">
									<div class="text-xs">{t.label}</div>
									<div class="text-xs text-slate-400 mt-0.5">{t.desc}</div>
								</div>
							</button>
						{/each}
					</div>
				</div>
				<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6">
					<h3 class="text-xs font-semibold text-white flex items-center gap-2 mb-4">
						<Info size={13} class="text-slate-400" /> Environment Info
					</h3>
					<div class="flex flex-col gap-3 text-xs font-mono">
						<div class="flex justify-between border-b border-white/[0.02] pb-2.5">
							<span class="text-slate-300">Wails Version</span>
							<strong class="text-slate-300">v2.12.0</strong>
						</div>
						<div class="flex justify-between border-b border-white/[0.02] pb-2.5">
							<span class="text-slate-300">Frontend Engine</span>
							<strong class="text-slate-300">Svelte 5 + Vite</strong>
						</div>
						<div class="flex justify-between border-b border-white/[0.02] pb-2.5">
							<span class="text-slate-300">Backend Language</span>
							<strong class="text-slate-300">Go (1.26.3)</strong>
						</div>
						<div class="flex justify-between pb-1">
							<span class="text-slate-300">Runtime Env</span>
							<strong class="text-emerald-500 flex items-center gap-1">
								<ShieldCheck size={12} /> OS Native
							</strong>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
