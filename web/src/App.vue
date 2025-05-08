<script>
import Navbar from '@/components/Navbar.vue';
import { useAuthStore } from "@/stores/authStore";

export default {
  components: { Navbar },
  async mounted() {
    const auth = useAuthStore();
    await auth.checkLogin();
  },
};
</script>

<template>
  <div class="bg-yellow-50 text-foreground_primary min-h-screen bg-[url('/OldPaperTexture.png')] bg-repeat">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" />
    <Navbar />
    <router-view />
  </div>
</template>
