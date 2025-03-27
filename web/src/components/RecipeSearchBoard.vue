<script setup>

import {ref, onMounted, defineProps, watch} from 'vue';
import Card from '@/components/ViewRecipeCard.vue';

const recipes = ref([]);

const props = defineProps({
  searchTerm: String,
});

const fetchRecipes = async () => {

  try {

    const response = await fetch (`/api/data/recipe/search?searchTerm=${encodeURIComponent(props.searchTerm)}`, {
      method: 'GET',  
      headers: { 'Content-Type': 'application/json' },
    });

    if (!response.ok) throw new Error('Failed to fetch');

    const data = await response.json();
    recipes.value = data;
    
  } catch (error) {

    alert('Failed to fetch recipes casue, ' + error); 
    
  }

}; 

watch(() => props.searchTerm, fetchRecipes);

onMounted(fetchRecipes);

</script>

<template>

  <div class="grid md:grid-cols-3 place-items-center gap-y-12 pb-12">
    <div v-for="recipe in recipes">
      <Card :name="recipe.name" :views="recipe.views" :uuid="recipe.uuid" />
    </div>
  </div>  

</template>
