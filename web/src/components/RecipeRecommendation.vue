<script setup>

import {ref, onMounted} from 'vue';
import Card from '@/components/ViewRecipeCard.vue';

const recipes = ref([]);

const mostViewed = async () => {

  try {

    const response = await fetch ('/api/data/recipe/mostviewed', {
      method: 'GET',  
    });

    if (!response.ok) throw new Error('Failed to fetch');

    const data = await response.json();
    recipes.value = data;

    
  } catch (error) {

    alert('Failed to fetch recipes casue, ' + error); 
    
  }

}; 

onMounted(mostViewed);

</script>

<template>

  <h2 class="text-center text-amber-900 font-[Monomakh] font-bold text-3xl pt-12 pb-12">Most Viewed</h2>

  <div class="grid md:grid-cols-3 place-items-center gap-0 pb-12">
    <div v-for="recipe in recipes">
      <Card :name="recipe.name" :views="recipe.views" :uuid="recipe.uuid" />
    </div>
  </div>  

</template>
