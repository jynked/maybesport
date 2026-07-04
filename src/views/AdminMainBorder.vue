<template>
    <main>
        <div class="admin-main-border">
            <h1>Редактирование главного баннера</h1>

            <div v-if="loading" class="loading">Загрузка...</div>
            <div v-else class="form-container">
                <div class="form-group">
                    <label>Заголовок (RU)</label>
                    <input type="text" v-model="form.title_ru" />
                </div>
                <div class="form-group">
                    <label>Заголовок (EN)</label>
                    <input type="text" v-model="form.title_en" />
                </div>

                <div class="form-group">
                    <label>ID товара для редиректа (uniqueId, например 1-1)</label>
                    <input type="text" v-model="form.newItemUniqueId" />
                </div>

                <div class="form-group">
                    <label>Медиа (изображение / GIF / видео → GIF)</label>
                    <input type="file" accept="image/jpeg,image/png,image/gif,video/mp4,video/webm,video/quicktime"
                        @change="onFileSelected" :disabled="uploadingFile" />
                    <div v-if="uploadingFile" class="uploading">Загрузка файла...</div>
                    <div v-if="mediaPreview" class="media-preview">
                        <img v-if="isImagePreview" :src="mediaPreview" alt="Preview" />
                        <video v-else-if="isVideoPreview" :src="mediaPreview" controls muted loop
                            style="max-width: 300px;"></video>
                        <div v-else class="file-name">Выбран: {{ selectedFileName }}</div>
                    </div>
                    <div v-if="form.image && !mediaPreview" class="image-preview">
                        <img :src="form.image" alt="Current banner" />
                        <button class="remove-preview" @click="clearMedia">✖</button>
                    </div>
                </div>

                <button class="save-btn" @click="save" :disabled="saving || uploadingFile">
                    <span v-if="saving">Сохранение...</span>
                    <span v-else>Сохранить</span>
                </button>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../api';
import { useToastStore } from '../stores/toast';

const toast = useToastStore();
const loading = ref(false);
const saving = ref(false);
const uploadingFile = ref(false);
const selectedFile = ref(null);
const mediaPreview = ref(null);
const selectedFileName = ref('');

const emit = defineEmits(['page-loaded']);

const form = ref({
    title_ru: '',
    title_en: '',
    image: '',
    newItemUniqueId: ''
});

const isImagePreview = computed(() => {
    if (!mediaPreview.value) return false;
    return mediaPreview.value.startsWith('data:image') || /\.(jpg|jpeg|png|gif)$/i.test(mediaPreview.value);
});

const isVideoPreview = computed(() => {
    if (!mediaPreview.value) return false;
    return mediaPreview.value.startsWith('data:video') || /\.(mp4|webm|mov)$/i.test(mediaPreview.value);
});

async function loadData() {
    loading.value = true;
    try {
        const resp = await api.getMainPageNew();
        const data = resp.data;
        form.value.title_ru = data.title?.ru || '';
        form.value.title_en = data.title?.en || '';
        form.value.image = data.image || '';
        form.value.newItemUniqueId = data.newItemUniqueId || '1-1';
        emit('page-loaded', true)
    } catch (err) {
        console.error(err);
        toast.error('Не удалось загрузить данные баннера');
    } finally {
        loading.value = false;
    }
}

function onFileSelected(event) {
    const file = event.target.files[0];
    if (!file) return;
    selectedFile.value = file;
    selectedFileName.value = file.name;

    const reader = new FileReader();
    reader.onload = (e) => {
        mediaPreview.value = e.target.result;
    };
    reader.readAsDataURL(file);
}

async function uploadCurrentFile() {
    if (!selectedFile.value) return null;
    uploadingFile.value = true;
    try {
        const resp = await api.uploadMedia(selectedFile.value);
        return resp.data.url;
    } catch (err) {
        console.error(err);
        toast.error('Ошибка загрузки медиафайла');
        return null;
    } finally {
        uploadingFile.value = false;
    }
}

async function save() {
    let finalImageUrl = form.value.image;

    if (selectedFile.value) {
        const uploadedUrl = await uploadCurrentFile();
        if (!uploadedUrl) return;
        finalImageUrl = uploadedUrl;
    }

    if (!finalImageUrl || !finalImageUrl.trim()) {
        toast.error('Медиафайл обязателен');
        return;
    }

    saving.value = true;
    try {
        await api.updateMainPage({
            title_ru: form.value.title_ru,
            title_en: form.value.title_en,
            image: finalImageUrl,
            newItemUniqueId: form.value.newItemUniqueId
        });
        toast.success('Баннер успешно обновлён');
        form.value.image = finalImageUrl;
        selectedFile.value = null;
        mediaPreview.value = null;
        selectedFileName.value = '';
    } catch (err) {
        console.error(err);
        toast.error('Ошибка при сохранении');
    } finally {
        saving.value = false;
    }
}

function clearMedia() {
    form.value.image = '';
    selectedFile.value = null;
    mediaPreview.value = null;
    selectedFileName.value = '';
}

onMounted(() => {
    loadData();
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/admin-main-border';
</style>