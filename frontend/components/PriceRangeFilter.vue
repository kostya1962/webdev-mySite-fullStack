<template>
<div class="price-filter">
    <div class="slider-wrap">
    <div class="track" :style="trackStyle" aria-hidden="true"></div>

    <input
        class="range range-min"
        type="range"
        :min="min"
        :max="max"
        :step="1"
        v-model.number="localMin"
        @input="onInputMin"
        @change="emitPrice"
        :aria-label="`Минимальная цена`"
    />

    <input
        class="range range-max"
        type="range"
        :min="min"
        :max="max"
        :step="1"
        v-model.number="localMax"
        @input="onInputMax"
        @change="emitPrice"
        :aria-label="`Максимальная цена`"
    />
    </div>

    <div class="values">Цена: {{ localMin }} ₽ – {{ localMax }} ₽</div>
</div>
</template>

<script setup lang="ts">
import { ref, computed, watch, type Ref } from 'vue'

type PriceTuple = [number, number]

const props = defineProps<{
price?: PriceTuple
min?: number
max?: number
}>()

const emit = defineEmits<{ (e: 'update:price', value: PriceTuple): void }>()

const min = props.min ?? 40
const max = props.max ?? 180

const initial: PriceTuple = [
Math.max(min, props.price?.[0] ?? min),
Math.min(max, props.price?.[1] ?? max)
]

const localMin: Ref<number> = ref(initial[0])
const localMax: Ref<number> = ref(initial[1])

watch(
() => props.price,
(val) => {
    if (Array.isArray(val) && val.length >= 2) {
    localMin.value = Math.max(min, Number(val[0]))
    localMax.value = Math.min(max, Number(val[1]))
    }
}
)

function onInputMin() {
if (localMin.value > localMax.value) {
    localMin.value = localMax.value
}
emitPrice()
}

function onInputMax() {
if (localMax.value < localMin.value) {
    localMax.value = localMin.value
}
emitPrice()
}

function emitPrice() {
emit('update:price', [Number(localMin.value), Number(localMax.value)])
}

const trackStyle = computed(() => {
const range = max - min || 1
const leftPct = ((localMin.value - min) / range) * 100
const rightPct = ((localMax.value - min) / range) * 100
return {
    background: `linear-gradient(90deg, var(--track-bg) ${leftPct}%, var(--accent) ${leftPct}%, var(--accent) ${rightPct}%, var(--track-bg) ${rightPct}%)`
}
})
</script>

<style scoped>
:root {
--track-height: 8px;
--thumb-size: 18px;
--track-bg: #e6e6e6;
--accent: #1167d1;
}

.price-filter {
width: 100%;
max-width: 420px;
font-family: system-ui, -apple-system, 'Segoe UI', Roboto, 'Helvetica Neue', Arial;
font-size: 14px;
color: #111;
}

.slider-wrap {
position: relative;
height: 36px;
display: flex;
align-items: center;
}

.track {
position: absolute;
left: 0;
right: 0;
height: var(--track-height);
border-radius: 999px;
top: 50%;
transform: translateY(-50%);
z-index: 1;
pointer-events: none;
}

.range {
-webkit-appearance: none;
appearance: none;
position: absolute;
left: 0;
right: 0;
width: 100%;
background: transparent;
pointer-events: auto;
}

.range::-webkit-slider-runnable-track {
height: var(--track-height);
background: transparent;
border: none;
}

.range::-moz-range-track {
height: var(--track-height);
background: transparent;
border: none;
}

.range::-webkit-slider-thumb {
-webkit-appearance: none;
appearance: none;
width: var(--thumb-size);
height: var(--thumb-size);
border-radius: 50%;
background: #fff;
border: 2px solid var(--accent);
box-shadow: 0 1px 2px rgba(0,0,0,0.12);
margin-top: calc((var(--track-height) - var(--thumb-size)) / 2);
pointer-events: auto;
cursor: pointer;
}

.range::-moz-range-thumb {
width: var(--thumb-size);
height: var(--thumb-size);
border-radius: 50%;
background: #fff;
border: 2px solid var(--accent);
box-shadow: 0 1px 2px rgba(0,0,0,0.12);
pointer-events: auto;
cursor: pointer;
}

.range:focus {
outline: none;
}

.range.range-min {
z-index: 3;
}

.range.range-max {
z-index: 4;
}

.values {
margin-top: 10px;
color: #222;
font-weight: 500;
}

@media (max-width: 420px) {
.price-filter { max-width: 100%; }
}
</style>
