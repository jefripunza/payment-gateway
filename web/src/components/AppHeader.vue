<script setup lang="ts">
import { Menu } from 'lucide-vue-next'

const drawerOpen = defineModel<boolean>('open', { default: false })

const props = defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <header
    class="sticky top-0 z-30 flex h-16 items-center gap-3 border-b border-white/5 bg-slate-950/70 px-4 backdrop-blur-xl sm:px-6"
  >
    <button
      type="button"
      class="rounded-lg p-2 text-slate-400 transition hover:bg-white/5 hover:text-slate-200 lg:hidden"
      @click="drawerOpen = true"
    >
      <Menu class="size-5" />
    </button>
    <div class="min-w-0">
      <h1 class="truncate text-lg font-bold tracking-tight text-slate-100">{{ props.title }}</h1>
      <p v-if="props.subtitle" class="hidden truncate text-xs text-slate-500 sm:block">
        {{ props.subtitle }}
      </p>
    </div>
    <div class="ml-auto flex items-center gap-2">
      <slot name="actions" />
    </div>
  </header>
</template>
