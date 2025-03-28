<script setup>

import {ref, onMounted, defineProps, watch} from 'vue';
import Card from '@/components/UserRecipeCard.vue';

const recipes = ref([]);

const fetchRecipes = async () => {

  try {

    const response = await fetch('/api/data/login/recipes', {
      credentials: "include",
      method: "GET",
    });

    if (!response.ok) throw new Error('Failed to fetch');

    const data = await response.json();
    recipes.value = data;
    console.log(recipes)
    
  } catch (error) {

    alert('Failed to fetch recipes casue, ' + error); 
    
  }

}; 

watch( () => recipes.value, (newValue, oldValue) => {
  console.log(`Views changed from ${oldValue} to ${newValue}`);
});
onMounted(fetchRecipes);

</script>

<template>

  <div class="grid md:grid-cols-3 place-items-center gap-y-12 pb-12">
    <div v-for="recipe in recipes">
      <Card :name="recipe.name" :views="recipe.views" :uuid="recipe.uuid" />
    </div>
  </div>  

</template>
