<script setup>
import { useRoute } from 'vue-router';
import { watch, ref, onMounted } from 'vue';

const route = useRoute();
const recipeId = ref(route.params.id);
const recipe = ref(null);
const getRecipe = async () => {
  if (!recipeId.value) return;

  try {
    const response = await fetch('/api/recipe/' + recipeId.value, {
      method: 'GET',
    });

    if (response.ok) {
      recipe.value = await response.text();
      console.log('Got the recipe', recipe.value);
    } else {
      console.error('Didn\'t get the recipe');
    }
  } catch (error) {
    console.error('Cannot get the recipe:', error);
  }
};
onMounted(getRecipe);
watch(() => route.params.id, (newId) => {
  recipeId.value = newId;
  getRecipe();
});
</script>
<template>
  <div class="output prose ml-50 my-20 
           lg:prose-li:marker:text-amber-900 lg:prose-headings:text-amber-900  
           lg:prose-headings:font-[Monomakh] lg:text-amber-900" v-html="recipe">
  </div>
</template>
