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
    <div class="absolute inset-0 bg-midnight-ink/20 backdrop-blur-sm" @click="open = false" />
    <div class="card-surface relative w-full max-w-md p-6 shadow-xl">
      <div class="flex items-start gap-3">
        <div class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-tangerine/10">
          <Trash2 class="size-5 text-tangerine" />
        </div>
        <div>
          <h3 class="text-base font-semibold text-midnight-ink">{{ title }}</h3>
          <p class="mt-1 text-xs text-fog">{{ description }}</p>
        </div>
      </div>

      <div class="mt-5">
        <label class="mb-1.5 block text-xs font-semibold text-graphite">
          Type <span class="font-mono text-tangerine">{{ confirmText }}</span> to confirm
        </label>
        <input
          v-model="typed"
          type="text"
          autocomplete="off"
          class="input-surface font-mono"
        />
        <p class="mt-1.5 text-[10px] text-silver">{{ itemName }} name must match exactly</p>
      </div>

      <div v-if="error" class="mt-3 rounded-lg border border-tangerine/20 bg-tangerine/5 px-3 py-2 text-xs text-tangerine">
        {{ error }}
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <button
          type="button"
          class="btn-ghost"
          @click="open = false"
        >
          Cancel
        </button>
        <button
          type="button"
          :disabled="!canDelete || saving"
          class="inline-flex items-center gap-2 rounded-lg bg-tangerine px-4 py-2 text-sm font-medium text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
          @click="emit('confirm')"
        >
          <Loader2 v-if="saving" class="size-3.5 animate-spin" />
          Delete
        </button>
      </div>
    </div>
  </div>
</template>
