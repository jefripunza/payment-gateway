<script setup lang="ts">
import { ref, watch } from 'vue'
import { Loader2, Trash2 } from 'lucide-vue-next'

const props = defineProps<{
  title: string
  description: string
  itemName: string
  confirmText: string
  error?: string
  saving?: boolean
}>()

const emit = defineEmits<{ confirm: [] }>()

const open = defineModel<boolean>('open', { default: false })

const typed = ref('')

watch(open, (v) => {
  if (v) typed.value = ''
})

const canDelete = ref(false)

watch(typed, (v) => {
  canDelete.value = v === props.confirmText
})
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-slate-950/70 backdrop-blur-sm" @click="open = false" />
    <div class="relative w-full max-w-md rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-2xl">
      <div class="flex items-start gap-3">
        <div class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-red-500/15">
          <Trash2 class="size-5 text-red-400" />
        </div>
        <div>
          <h3 class="text-base font-bold text-slate-100">{{ title }}</h3>
          <p class="mt-1 text-xs text-slate-400">{{ description }}</p>
        </div>
      </div>

      <div class="mt-5">
        <label class="mb-1.5 block text-xs font-medium text-slate-300">
          Type <span class="font-mono text-red-400">{{ confirmText }}</span> to confirm
        </label>
        <input
          v-model="typed"
          type="text"
          autocomplete="off"
          class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 font-mono text-sm text-slate-100 focus:border-red-400/60 focus:ring-2 focus:ring-red-500/30 focus:outline-none"
        />
        <p class="mt-1.5 text-[10px] text-slate-500">{{ itemName }} name must match exactly</p>
      </div>

      <div
        v-if="error"
        class="mt-3 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-300"
      >
        {{ error }}
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <button
          type="button"
          class="h-9 rounded-lg px-4 text-sm font-medium text-slate-400 transition hover:bg-white/5 hover:text-slate-200"
          @click="open = false"
        >
          Cancel
        </button>
        <button
          type="button"
          :disabled="!canDelete || saving"
          class="flex h-9 items-center gap-2 rounded-lg bg-red-600 px-4 text-sm font-semibold text-white transition hover:bg-red-500 disabled:cursor-not-allowed disabled:opacity-50"
          @click="emit('confirm')"
        >
          <Loader2 v-if="saving" class="size-3.5 animate-spin" />
          Delete
        </button>
      </div>
    </div>
  </div>
</template>
