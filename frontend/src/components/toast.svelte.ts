export interface Toast {
	id: string;
	type: 'success' | 'error' | 'info' | 'warning';
	message: string;
	duration?: number;
}

let toasts = $state<Toast[]>([]);

export const toastState = {
	get items() { return toasts; },
	add(type: Toast['type'], message: string, duration = 3000) {
		const id = Math.random().toString(36).substring(2, 9);
		toasts.push({ id, type, message, duration });
		setTimeout(() => this.remove(id), duration);
	},
	remove(id: string) {
		toasts = toasts.filter(t => t.id !== id);
	}
};
