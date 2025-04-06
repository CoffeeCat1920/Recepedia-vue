<script setup>

import { ref, onMounted, defineProps, watch } from 'vue';
import Card from '@/components/AdminUserCard.vue';

const users = ref([]);

const fetchUsers = async () => {

  try {
    const response = await fetch('/api/data/admin/allusers', {
      credentials: "include",
      method: "GET",
    });

    if (!response.ok) throw new Error('Failed to fetch');

    const data = await response.json();
    users.value = data;
  } catch (error) {
    alert('Failed to fetch users casue, ' + error);
  }

};

onMounted(fetchUsers);

</script>

<template>

  <div class="grid md:grid-cols-3 place-items-center gap-y-12 pb-12">
    <div v-for="user in users">
      <Card :name="user.name" :uuid="user.uuid" />
    </div>
  </div>

</template>
