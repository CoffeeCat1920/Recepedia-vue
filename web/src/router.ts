import { createRouter, createWebHistory } from 'vue-router';
import Home from './views/Home.vue';
import About from './views/About.vue';
import Signup from './views/Signup.vue';
import Login from './views/Login.vue';
import Dashboard from './views/Dashboard.vue';
import UploadRecipe from './views/UploadRecipe.vue';
import Recipe from './views/Recipe.vue'
import Search from './views/Search.vue'
import EditRecipe from './views/EditRecipe.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
  { path: '/signup', component: Signup },
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard },
  { path: '/uploadrecipe', component: UploadRecipe },
  { path: '/recipe/:id', component: Recipe},
  { path: '/search', component: Search},
  { path: '/editrecipe/:id', component: EditRecipe},
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
