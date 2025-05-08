<script setup>
import { useRoute } from 'vue-router';
import { watch, ref, onMounted } from 'vue';

const route = useRoute();
const recipeId = ref(route.params.id);
const recipe = ref(null);
const isLiked = ref(false);
const userId = 'some-user-id'; // Replace this with your actual user ID (auth context or prop)

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

const checkLiked = async () => {
  try {
    const response = await fetch('/api/recipe/isliked', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userid: userId, recipeid: recipeId.value })
    });

    if (response.ok) {
      const data = await response.json();
      isLiked.value = data.liked;
    } else {
      isLiked.value = false;
    }
  } catch (error) {
    console.error('Error checking like:', error);
  }
};

const toggleLike = async () => {
  const url = isLiked.value ? '/api/recipe/unlike' : '/api/recipe/like';

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userid: userId, recipeid: recipeId.value })
    });

    if (response.ok) {
      isLiked.value = !isLiked.value;
    } else {
      console.error('Failed to toggle like');
    }
  } catch (error) {
    console.error('Error toggling like:', error);
  }
};

onMounted(() => {
  getRecipe();
  checkLiked();
});

watch(() => route.params.id, (newId) => {
  recipeId.value = newId;
  getRecipe();
  checkLiked();
});
</script>

<template>
  <div class="relative">
    <!-- Like button -->
    <button @click="toggleLike" class="absolute top-0 right-0 m-4 text-3xl cursor-pointer transition-all duration-200"
      :aria-label="isLiked ? 'Unlike' : 'Like'">
      <span :class="{ 'text-red-500': isLiked, 'text-gray-400': !isLiked }">
        {{ isLiked ? '❤️' : '🤍' }}
      </span>
    </button>

    <!-- Recipe content -->
    <div class="output prose ml-50 my-20
             lg:prose-li:marker:text-amber-900 lg:prose-headings:text-amber-900  
             lg:prose-headings:font-[Monomakh] lg:text-amber-900" v-html="recipe" />
  </div>
</template>
