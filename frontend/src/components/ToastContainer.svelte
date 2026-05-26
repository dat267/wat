<script lang="ts">
	import { CheckCircle, XCircle, Info, AlertTriangle, X } from '@lucide/svelte';
	import { toastState } from './toast.svelte';
	import { fade, fly } from 'svelte/transition';
	const icons = {
		success: CheckCircle,
		error: XCircle,
		info: Info,
		warning: AlertTriangle
	};
	const colors = {
		success: 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400',
		error: 'bg-rose-500/10 border-rose-500/20 text-rose-400',
		info: 'bg-blue-500/10 border-blue-500/20 text-blue-400',
		warning: 'bg-amber-500/10 border-amber-500/20 text-amber-400'
	};
</script>
<div class="fixed top-6 right-6 z-[999] flex flex-col gap-3 pointer-events-none select-none max-w-sm w-full">
	{#each toastState.items as toast (toast.id)}
		{@const Icon = icons[toast.type]}
		<div
			in:fly={{ y: -20, duration: 300 }}
			out:fade={{ duration: 150 }}
			class="flex items-start gap-3 p-4 rounded-2xl border backdrop-blur-2xl shadow-2xl pointer-events-auto transition-all {colors[toast.type]}"
		>
			<div class="shrink-0 mt-0.5">
				<Icon size={16} />
			</div>
			<div class="flex-1 text-xs font-semibold leading-relaxed font-sans">{toast.message}</div>
			<button
				onclick={() => toastState.remove(toast.id)}
				class="shrink-0 p-0.5 -mt-1 -mr-1 rounded-lg hover:bg-white/[0.05] transition-colors cursor-pointer text-slate-400 hover:text-white"
			>
				<X size={12} />
			</button>
		</div>
	{/each}
</div>
