<template>
    <main>
        <div class="container-fluid py-4">
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h1>Управление пользователями</h1>
            </div>

            <div class="card mb-4">
                <div class="card-body">
                    <div class="row">
                        <div class="col-md-6">
                            <input type="text" class="form-control" placeholder="Поиск по email или имени"
                                v-model="searchQuery" @input="onSearchInput" />
                        </div>
                        <div class="col-md-2">
                            <button class="btn btn-primary w-100" @click="loadUsers">Поиск</button>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="loading" class="text-center py-5">
                <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Загрузка...</span>
                </div>
            </div>
            <div v-else-if="users.length === 0" class="alert alert-info">Пользователи не найдены</div>
            <div v-else class="table-responsive">
                <table class="table table-striped table-hover">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Email</th>
                            <th>Имя</th>
                            <th>Дата регистрации</th>
                            <th>Статус</th>
                            <th>Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in users" :key="user.id">
                            <td>{{ user.id }}</td>
                            <td>{{ user.email }}</td>
                            <td>{{ user.name || '—' }}</td>
                            <td>{{ formatDate(user.created_at) }}</td>
                            <td>
                                <span :class="user.is_banned ? 'badge bg-danger' : 'badge bg-success'">
                                    {{ user.is_banned ? 'Заблокирован' : 'Активен' }}
                                </span>
                            </td>
                            <td>
                                <button class="btn btn-sm"
                                    :class="user.is_banned ? 'btn-outline-success' : 'btn-outline-danger'"
                                    @click="toggleBan(user)" :disabled="actionLoading">
                                    {{ user.is_banned ? 'Разблокировать' : 'Заблокировать' }}
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div class="pagination" v-if="totalPages > 1">
                <button :disabled="page === 1" @click="changePage(page - 1)">«</button>
                <span>Страница {{ page }} из {{ totalPages }}</span>
                <button :disabled="page === totalPages" @click="changePage(page + 1)">»</button>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { http } from '../api';
import { useToastStore } from '../stores/toast';
import { useAuthStore } from '../stores/auth';

const emit = defineEmits(['page-loaded']);

const users = ref([]);
const loading = ref(false);
const actionLoading = ref(false);
const page = ref(1);
const totalPages = ref(1);
const searchQuery = ref('');
const limit = 20;

let searchTimeout = null;

function formatDate(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    return d.toLocaleString();
}

async function loadUsers() {
    loading.value = true;
    try {
        const params = {
            page: page.value,
            limit: limit,
            search: searchQuery.value || undefined
        };
        const response = await http.get('/admin/users', { params });
        users.value = response.data.users;
        totalPages.value = response.data.totalPages;
    } catch (error) {
        console.error(error);
        useToastStore().error('Ошибка загрузки пользователей');
    } finally {
        loading.value = false;
        emit('page-loaded', true);
    }
}

function changePage(newPage) {
    page.value = newPage;
    loadUsers();
}

function onSearchInput() {
    if (searchTimeout) clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
        page.value = 1;
        loadUsers();
    }, 500);
}

async function toggleBan(user) {
    if (actionLoading.value) return;
    actionLoading.value = true;
    try {
        await http.put(`/admin/users/${user.id}/ban`, { ban: !user.is_banned });
        useToastStore().success(user.is_banned ? 'Пользователь разблокирован' : 'Пользователь заблокирован');
        await loadUsers();
    } catch (error) {
        console.error(error);
        useToastStore().error('Ошибка изменения статуса');
    } finally {
        actionLoading.value = false;
    }
}

onMounted(() => {
    loadUsers();
});
</script>

<style scoped>
.pagination {
    margin-top: 20px;
    display: flex;
    justify-content: center;
    gap: 10px;
    align-items: center;
}

.pagination button {
    padding: 5px 10px;
    border: 1px solid #dee2e6;
    background: white;
    cursor: pointer;
}

.pagination button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
</style>