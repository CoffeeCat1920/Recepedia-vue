!-- /component/Navbar.vue -->

<script>
  import { useAuthStore } from "@/stores/authStore";
  import { useRouter } from "vue-router";

  export default {
    setup() {
      const router = useRouter();
      const auth = useAuthStore();

      const logout = async () => {
        await auth.logout(router);
      };

      return { auth, logout };
    },
    async mounted() {
      await this.auth.checkLogin(); 
    },
  };
</script>

<template>

  <nav class="navbar">

    <router-link class="brand" to="/">•ᴗ• Recipidia</router-link>

    <div class="nav-links">

      <ul class="nav-list">

        <li class="left-just">
          <router-link to="/" class="nav-item">
            Home
          </router-link>
        </li>

        <li class="left-just">
            <router-link to="/about" class="nav-item"> 
            About 
          </router-link>
        </li>

        <template v-if="auth.loggedIn">
          <li class="left-just" ><router-link to="/uploadrecipe" class="nav-item" >Upload Recipe</router-link></li>
          <li class="left-just" ><router-link to="/dashboard" class="nav-item" >Dashboard</router-link></li>
          <li class="left-just" >
            <a href="#" @click.prevent="logout" class="nav-item">
              Logout
            </a>
          </li>
        </template>

        <template v-else>
          <li class="left-just"><router-link to="/signup" class="nav-item">SignUp</router-link></li>

          <li class="left-just"><router-link to="/login" class="nav-item">LogIn </router-link></li>
        </template>


      </ul>

    </div>

  </nav>

</template>
