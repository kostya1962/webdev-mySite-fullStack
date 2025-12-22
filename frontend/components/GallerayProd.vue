<script setup lang="ts">
    import { ref, computed, type CSSProperties } from 'vue'

    const API_URLimage = useAPIimage();

    const props = defineProps<{
    images: string[]
    }>()

    const activeIndex = ref<number>(0)

    const setActive = (index: number): void => {
        activeIndex.value = index
    }

    const progressStyle = computed<CSSProperties>(() => {
        const count = props.images.length || 1
        const width = 100 / count

        return {
            width: `${width}%`,
            transform: `translateX(${activeIndex.value * 100}%)`
        }
    });

    function next() { 
        activeIndex.value = (activeIndex.value + 1) % props.images.length 
    }
</script>

<template>
    <div class="product-gallery">

        <div class="thumbnails">
        <img
            v-for="(img, index) in images"
            :key="img"
            :src="API_URLimage + img"
            class="thumb"
            :class="{ active: index === activeIndex }"
            @click="setActive(index)"
        />
        </div>


        <div class="preview" :style="{ backgroundImage: `url(${API_URLimage}${images[activeIndex]})` } " @click="next">
            <div class="progress">
                <span
                    class="progress__active"
                    :style="progressStyle"
                />
            </div>
        </div>
    </div>
</template>

<style scoped>
    .product-gallery {
    display: flex;
    gap: 32px;
    max-width: 700px;
    }

    .thumbnails {
    display: flex;
    flex-direction: column;
    gap: 16px;
    }

    .thumb {
    width: 60px;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 0.3s ease;
    }

    .thumb.active {
    opacity: 1;
    }

    .preview {
        aspect-ratio: 1/1;
        border-radius: 8px;
        min-width: 320px;
        width: 100%;    
        flex: 1;
        background-size: cover;
        background-position: center;
        background-repeat: no-repeat;
        background-color: lightgray;
        display: flex;
        gap:16px;
        flex-direction: column;
        padding-bottom: -6px; 
    }

    .preview__image {
    width: 100%;
    display: block;
    }

    .progress {
        position: relative;
        height: 2px;
        background: #e5e5e5;
        margin-top: 16px;
        overflow: hidden;
        margin-top: auto
        }

    .progress__active {
    position: absolute;
    top: 0;
    left: 0;
    height: 2px;
    background: #000;
    transition: transform 0.3s ease;
    }
</style>