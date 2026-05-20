<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-5" :class="bgClass">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">
          {{ label }}
        </p>
        <p class="mt-1.5 font-bold text-gray-900 truncate" :class="small ? 'text-lg' : 'text-2xl'">
          {{ value }}
        </p>
      </div>
      <span class="text-2xl">{{ icon }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  icon: { type: String, default: '📊' },
  label: { type: String, required: true },
  value: { type: [String, Number], default: '—' },
  color: { type: String, default: 'blue' },
  small: { type: Boolean, default: false },
})

const bgClass = computed(() => ({
  'border-l-4 border-l-blue-400': props.color === 'blue',
  'border-l-4 border-l-purple-400': props.color === 'purple',
  'border-l-4 border-l-yellow-400': props.color === 'yellow',
  'border-l-4 border-l-green-400': props.color === 'green',
}))
</script>
