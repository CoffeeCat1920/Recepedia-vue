import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from "@/stores/authStore";

import Home from './views/Home.vue';
import About from './views/About.vue';
import Signup from './views/Signup.vue';
import Login from './views/Login.vue';
import Dashboard from './views/Dashboard.vue';
import UploadRecipe from './views/UploadRecipe.vue';
import Recipe from './views/Recipe.vue'
import Search from './views/Search.vue'
import EditRecipe from './views/EditRecipe.vue'
import AdminDashboard from './views/AdminDashboard.vue'

import AdminEditRecipes from './views/AdminEditRecipes.vue'
import AdminEditUsers from './views/AdminEditUsers.vue'

import AdminLogin from './views/AdminLogin.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
  { path: '/signup', component: Signup },
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard },
  { path: '/uploadrecipe', component: UploadRecipe },
  { path: '/recipe/:id', component: Recipe },
  { path: '/search', component: Search },
  { path: '/editrecipe/:id', component: EditRecipe },

  { path: '/adminlogin', component: AdminLogin },
  { path: '/admindashboard', component: AdminDashboard },

  { path: '/admin/editRecipes', component: AdminEditRecipes},
  { path: '/admin/editUsers', component: AdminEditUsers},
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  if (to.path === '/login' && authStore.loggedIn) {
    next('/dashboard');
  } else if (to.path === '/adminlogin' && authStore.adminLoggedIn) {
    next('/admindashboard');
  } else {
    next();
  }

});

export default router;
