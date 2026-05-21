<template>
  <div
    class="rounded-full overflow-hidden bg-white border border-gray-100 flex items-center justify-center shrink-0"
    :style="{ width: sizePx, height: sizePx }"
  >
    <img
      v-if="logoUrl"
      :src="logoUrl"
      :alt="name"
      class="w-full h-full object-contain p-0.5"
      @error="onError"
    />
    <span v-else class="text-sm">{{ fallbackEmoji }}</span>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  name: { type: String, required: true },
  size: { type: Number, default: 28 },
})

const LOGOS = {
  'Manchester City': 'https://upload.wikimedia.org/wikipedia/en/e/eb/Manchester_City_FC_badge.svg',
  Liverpool: 'https://upload.wikimedia.org/wikipedia/en/0/0c/Liverpool_FC.svg',
  Arsenal: 'https://upload.wikimedia.org/wikipedia/en/5/53/Arsenal_FC.svg',
  Chelsea: 'https://upload.wikimedia.org/wikipedia/en/c/cc/Chelsea_FC.svg',
}

const EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

const hasError = ref(false)

const logoUrl = computed(() => (hasError.value ? null : LOGOS[props.name] || null))

const fallbackEmoji = computed(() => EMOJIS[props.name] || '⚽')

const sizePx = computed(() => `${props.size}px`)

function onError() {
  hasError.value = true
}
</script>
