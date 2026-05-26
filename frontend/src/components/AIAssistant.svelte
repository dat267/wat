<script lang="ts">
	import { Bot, Send, RefreshCw, AlertTriangle, Key, Sparkles } from '@lucide/svelte';
	import { AskAI, StorePut, StoreGet } from '../../wailsjs/go/main/App';
	import { toastState } from './toast.svelte';
	interface Message {
		role: 'user' | 'assistant';
		content: string;
		time: string;
	}
	let provider = $state<'openai' | 'anthropic'>('anthropic');
	let prompt = $state('');
	let loading = $state(false);
	let messages = $state<Message[]>([]);
	let chatContainer = $state<HTMLDivElement | null>(null);
	const systemPrompts = {
		general: "You are a professional, helpful local desktop systems assistant. Keep answers brief, compact, direct and precise.",
		diagnostics: "You are an expert system administration assistant. Analyze standard system metrics and logs to diagnose hardware or software bottlenecks."
	};
	$effect(() => {
		loadHistory();
	});
	$effect(() => {
		if (messages.length && chatContainer) {
			chatContainer.scrollTo({ top: chatContainer.scrollHeight, behavior: 'smooth' });
		}
	});
	async function loadHistory() {
		try {
			const saved = await StoreGet('ai_chat_history', 'latest');
			if (saved) {
				messages = JSON.parse(saved);
			}
		} catch (e) {
			console.error(e);
		}
	}
	async function saveHistory() {
		try {
			await StorePut('ai_chat_history', 'latest', JSON.stringify(messages));
		} catch (e) {
			console.error(e);
		}
	}
	async function sendMessage(text = prompt) {
		if (!text.trim() || loading) return;
		const userMsg: Message = {
			role: 'user',
			content: text.trim(),
			time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
		};
		messages.push(userMsg);
		if (text === prompt) prompt = '';
		loading = true;
		try {
			const activeSysPrompt = text.toLowerCase().includes('log') || text.toLowerCase().includes('diagnose') ? systemPrompts.diagnostics : systemPrompts.general;
			const res = await AskAI(provider, '', activeSysPrompt, userMsg.content);
			messages.push({
				role: 'assistant',
				content: res,
				time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
			});
			saveHistory();
		} catch (err: any) {
			toastState.add('error', err?.toString() ?? 'Failed to communicate with AI endpoint.');
		} finally {
			loading = false;
		}
	}
	function clearChat() {
		messages = [];
		saveHistory();
		toastState.add('info', 'Chat history cleared.');
	}
</script>
<div class="space-y-8 select-none max-w-5xl h-[calc(100vh-4rem)] flex flex-col justify-between">
	<div>
		<h1 class="text-2xl font-bold tracking-tight text-white mb-1">AI Assistant</h1>
		<p class="text-slate-400 text-xs">Troubleshoot, write automation scripts, and audit diagnostics with AI.</p>
	</div>
	<div class="flex-1 bg-white/[0.01] border border-white/[0.03] rounded-2xl p-6 my-4 flex flex-col justify-between overflow-hidden relative backdrop-blur-3xl">
		<div class="flex items-center justify-between border-b border-white/[0.03] pb-4 mb-4">
			<div class="flex items-center gap-3">
				<div class="p-2 rounded-xl bg-white/[0.03] border border-white/[0.05] text-blue-400">
					<Bot size={16} />
				</div>
				<div>
					<h3 class="text-xs font-semibold text-white">Interactive Copilot</h3>
					<span class="text-[9px] font-mono text-slate-500 uppercase tracking-widest">Active Model</span>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<select
					bind:value={provider}
					class="bg-[#121726]/80 border border-white/[0.04] rounded-xl px-3 py-1.5 text-[11px] text-slate-300 focus:outline-none"
				>
					<option value="anthropic">Claude 3.5 Haiku</option>
					<option value="openai">GPT-4o Mini</option>
				</select>
				{#if messages.length > 0}
					<button
						onclick={clearChat}
						class="px-3 py-1.5 rounded-xl border border-rose-500/20 text-rose-400 text-[10px] font-mono hover:bg-rose-500/5 transition-colors cursor-pointer"
					>
						CLEAR
					</button>
				{/if}
			</div>
		</div>
		<div bind:this={chatContainer} class="flex-1 overflow-y-auto pr-2 space-y-4 flex flex-col justify-start">
			{#if messages.length === 0}
				<div class="my-auto text-center space-y-6 max-w-sm mx-auto">
					<div class="inline-flex p-4 rounded-3xl bg-blue-500/5 text-blue-400 animate-pulse border border-blue-500/10 shadow-2xl">
						<Sparkles size={32} />
					</div>
					<div class="space-y-2">
						<h4 class="text-xs font-bold text-white font-mono tracking-widest uppercase">Consult Diagnostics AI</h4>
						<p class="text-slate-400 text-xs leading-relaxed">
							Need a system automation script or diagnostic explanation? Ask me directly.
						</p>
					</div>
					<div class="grid grid-cols-2 gap-2.5 pt-2">
						<button
							onclick={() => sendMessage('Analyze my system resource metrics.')}
							class="p-2.5 rounded-xl border border-white/[0.03] bg-white/[0.005] hover:bg-white/[0.02] text-[10px] text-left text-slate-300 font-mono transition-all hover:scale-[1.02] cursor-pointer"
						>
							Analyze metrics
						</button>
						<button
							onclick={() => sendMessage('Write a shell script to automate docker cleanup.')}
							class="p-2.5 rounded-xl border border-white/[0.03] bg-white/[0.005] hover:bg-white/[0.02] text-[10px] text-left text-slate-300 font-mono transition-all hover:scale-[1.02] cursor-pointer"
						>
							Docker cleanup
						</button>
					</div>
				</div>
			{:else}
				{#each messages as msg}
					<div class="flex flex-col {msg.role === 'user' ? 'items-end' : 'items-start'} max-w-2xl {msg.role === 'user' ? 'ml-auto' : ''}">
						<div class="flex items-center gap-2 mb-1">
							<span class="text-[9px] font-bold font-mono tracking-widest uppercase text-slate-500">{msg.role}</span>
							<span class="text-[8px] font-mono text-slate-600">{msg.time}</span>
						</div>
						<div class="p-3.5 rounded-2xl border text-xs font-medium leading-relaxed font-mono whitespace-pre-wrap {msg.role === 'user' ? 'bg-[#1e293b]/20 border-white/[0.05] text-white' : 'bg-white/[0.01] border-white/[0.03] text-slate-300'}">
							{msg.content}
						</div>
					</div>
				{/each}
			{/if}
			{#if loading}
				<div class="flex items-center gap-2.5 text-slate-400 text-[10px] font-mono py-2 bg-white/[0.005] border border-white/[0.03] px-3.5 rounded-xl self-start">
					<RefreshCw size={11} class="animate-spin" /> Thinking…
				</div>
			{/if}
		</div>
		<div class="border-t border-white/[0.03] pt-4 mt-4 flex gap-2">
			<input
				type="text"
				bind:value={prompt}
				disabled={loading}
				placeholder="Ask for system script automations, logs analytics, diagnostic guides..."
				onkeydown={(e) => e.key === 'Enter' && sendMessage()}
				class="flex-1 bg-white/[0.01] border border-white/[0.04] rounded-xl px-4 py-3 text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-white/[0.08] font-mono transition-all"
			/>
			<button
				onclick={() => sendMessage()}
				disabled={loading || !prompt.trim()}
				class="px-5 rounded-xl bg-white text-slate-950 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-200 text-xs font-semibold flex items-center gap-1.5 transition-all active:scale-95 cursor-pointer font-mono"
			>
				<Send size={12} /> SEND
			</button>
		</div>
	</div>
</div>
