<script lang="ts">
	import { Terminal, Play, Plus, Trash2, Copy, CheckCircle, RefreshCw, Code } from '@lucide/svelte';
	import { ExecuteScript, StorePut, StoreGet } from '../../wailsjs/go/main/App';
	import { toastState } from './toast.svelte';
	interface CustomScript {
		name: string;
		command: string;
	}
	let activeScript = $state<'docker' | 'disk' | 'network' | 'services' | string>('disk');
	let output = $state('Click RUN script to start diagnostic output.');
	let loading = $state(false);
	let copied = $state(false);
	let customName = $state('');
	let customCmd = $state('');
	let customScriptsList = $state<CustomScript[]>([]);
	$effect(() => {
		loadCustomScripts();
	});
	async function loadCustomScripts() {
		try {
			const saved = await StoreGet('custom_scripts', 'list');
			if (saved) {
				customScriptsList = JSON.parse(saved);
			}
		} catch (e) {
			console.error(e);
		}
	}
	async function saveCustomScripts() {
		try {
			await StorePut('custom_scripts', 'list', JSON.stringify(customScriptsList));
		} catch (e) {
			console.error(e);
		}
	}
	async function runScript() {
		loading = true;
		output = 'Running diagnostic script execution, please wait...';
		try {
			let res = '';
			const isBuiltIn = ['docker', 'disk', 'network', 'services'].includes(activeScript);
			if (isBuiltIn) {
				res = await ExecuteScript(activeScript);
			} else {
				const sc = customScriptsList.find(s => s.name === activeScript);
				if (sc) {
					res = await ExecuteScript(sc.name);
				} else {
					throw new Error('Script not found');
				}
			}
			output = res || 'Script executed successfully but returned empty output.';
			toastState.add('success', 'Script executed successfully.');
		} catch (err: any) {
			output = err?.toString() ?? 'Error running script';
			toastState.add('error', 'Script execution failed.');
		} finally {
			loading = false;
		}
	}
	async function addCustomScript() {
		if (!customName.trim() || !customCmd.trim()) return;
		const name = customName.trim();
		const command = customCmd.trim();
		if (['docker', 'disk', 'network', 'services'].includes(name)) {
			toastState.add('warning', 'Name conflicts with built-in script.');
			return;
		}
		customScriptsList.push({ name, command });
		customName = '';
		customCmd = '';
		saveCustomScripts();
		toastState.add('success', 'Custom script registered.');
	}
	function removeCustomScript(name: string) {
		customScriptsList = customScriptsList.filter(s => s.name !== name);
		if (activeScript === name) activeScript = 'disk';
		saveCustomScripts();
		toastState.add('info', 'Custom script removed.');
	}
	function copyOutput() {
		if (!output) return;
		navigator.clipboard.writeText(output);
		copied = true;
		setTimeout(() => copied = false, 2000);
	}
</script>
<div class="space-y-8 select-none max-w-5xl">
	<div>
		<h1 class="text-2xl font-bold tracking-tight text-white mb-1">Script Console</h1>
		<p class="text-slate-400 text-xs">Run, manage, and register system shell script execution panels.</p>
	</div>
	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<div class="lg:col-span-2 flex flex-col gap-4">
			<div class="bg-[#090b11] border border-white/[0.04] rounded-2xl flex-1 flex flex-col justify-between overflow-hidden shadow-2xl relative">
				<div class="flex items-center justify-between px-5 py-3.5 border-b border-white/[0.03] bg-white/[0.005]">
					<div class="flex items-center gap-2.5 text-slate-300 text-xs font-semibold uppercase tracking-wider font-mono">
						<Terminal size={13} class="text-blue-400" />
						<span>Interactive Terminal</span>
					</div>
					<div class="flex items-center gap-3">
						<button
							onclick={copyOutput}
							disabled={loading}
							class="text-[10px] text-slate-400 hover:text-white flex items-center gap-1 cursor-pointer transition-colors font-mono disabled:opacity-40"
						>
							<Copy size={11} /> {copied ? 'COPIED!' : 'COPY'}
						</button>
						<button
							onclick={runScript}
							disabled={loading}
							class="px-3.5 py-1.5 rounded-lg bg-blue-500 hover:bg-blue-400 text-white text-[10px] font-bold flex items-center gap-1 transition-all active:scale-95 cursor-pointer font-mono"
						>
							{#if loading}
								<RefreshCw size={10} class="animate-spin" /> RUNNING...
							{:else}
								<Play size={10} /> RUN SCRIPT
							{/if}
						</button>
					</div>
				</div>
				<div class="flex-1 p-5 min-h-[300px] overflow-auto font-mono text-[11px] text-emerald-400 bg-black/40 whitespace-pre">
					{output}
				</div>
			</div>
		</div>
		<div class="space-y-6">
			<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-5 space-y-4">
				<h3 class="text-xs font-semibold text-white flex items-center justify-between">
					<span class="flex items-center gap-2">
						<Code size={13} class="text-slate-400" /> Select Script
					</span>
				</h3>
				<div class="flex flex-col gap-1.5 max-h-48 overflow-y-auto font-mono">
					{#each ['disk', 'network', 'docker', 'services'] as name}
						<button
							onclick={() => activeScript = name}
							class="w-full flex items-center justify-between text-left px-3.5 py-2.5 rounded-xl border text-[11px] transition-all cursor-pointer {activeScript === name ? 'border-white/[0.08] bg-white/[0.03] text-white font-semibold' : 'border-transparent bg-transparent text-slate-400 hover:bg-white/[0.01]'}"
						>
							<span>wat run {name}</span>
							<span class="text-[9px] uppercase tracking-widest text-slate-500 font-bold">Built-in</span>
						</button>
					{/each}
					{#each customScriptsList as s}
						<div class="w-full flex items-center gap-2">
							<button
								onclick={() => activeScript = s.name}
								class="flex-1 flex items-center justify-between text-left px-3.5 py-2.5 rounded-xl border text-[11px] transition-all cursor-pointer {activeScript === s.name ? 'border-white/[0.08] bg-white/[0.03] text-white font-semibold' : 'border-transparent bg-transparent text-slate-400 hover:bg-white/[0.01]'}"
							>
								<span class="truncate max-w-[120px]">{s.name}</span>
								<span class="text-[9px] uppercase tracking-widest text-blue-400 font-bold">Custom</span>
							</button>
							<button
								onclick={() => removeCustomScript(s.name)}
								class="p-2.5 rounded-xl border border-rose-500/20 text-rose-400 hover:bg-rose-500/5 transition-colors cursor-pointer"
							>
								<Trash2 size={12} />
							</button>
						</div>
					{/each}
				</div>
			</div>
			<div class="bg-white/[0.01] border border-white/[0.03] rounded-2xl p-5 space-y-4">
				<h3 class="text-xs font-semibold text-white flex items-center gap-2">
					<Plus size={13} class="text-slate-400" /> Custom Scripts
				</h3>
				<div class="space-y-3 font-mono">
					<div class="space-y-1">
						<label for="scripts-name-input" class="text-[9px] text-slate-400 uppercase tracking-wider">Command Name</label>
						<input
							id="scripts-name-input"
							type="text"
							bind:value={customName}
							placeholder="e.g. check-uptime"
							class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-3 py-2 text-xs text-white placeholder:text-slate-600 focus:outline-none"
						/>
					</div>
					<div class="space-y-1">
						<label for="scripts-cmd-input" class="text-[9px] text-slate-400 uppercase tracking-wider">Command Line</label>
						<input
							id="scripts-cmd-input"
							type="text"
							bind:value={customCmd}
							placeholder="e.g. uptime"
							class="w-full bg-white/[0.01] border border-white/[0.04] rounded-xl px-3 py-2 text-xs text-white placeholder:text-slate-600 focus:outline-none"
						/>
					</div>
					<button
						onclick={addCustomScript}
						disabled={!customName.trim() || !customCmd.trim()}
						class="w-full py-2.5 rounded-xl bg-white hover:bg-slate-200 disabled:opacity-40 text-slate-950 text-xs font-semibold flex items-center justify-center gap-1 transition-all cursor-pointer"
					>
						<Plus size={12} /> REGISTER SCRIPT
					</button>
				</div>
			</div>
		</div>
	</div>
</div>
