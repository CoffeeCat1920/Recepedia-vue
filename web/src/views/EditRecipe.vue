<script setup>
import { useRoute } from 'vue-router';
import { marked } from 'marked'
import { watch, computed, ref, onMounted } from 'vue';

const route = useRoute();
const recipeId = ref(route.params.id);
const recipeName = ref('//Name//');
const recipe = ref('# Please Write in md: [Guid](https://www.markdownguide.org/basic-syntax/)');

const getRecipe = async () => {
  if (!recipeId.value) return;

  try {
    const response = await fetch('/api/data/recipe/name/' + recipeId.value, {
      method: 'GET',
    });

    if (response.ok) {
      const data = await response.json();
      recipeName.value = data;
    } else {
      console.error('Didn\'t get the recipe');
    }
  } catch (error) {
    console.error('Cannot get the recipe:', error);
  }

  try {
    const response = await fetch('/api/data/recipe/content/' + recipeId.value, {
      method: 'GET',
    });

    if (response.ok) {
      const data = await response.text();
      recipe.value = data;
    } else {
      console.error('Didn\'t get the recipe');
    }
  } catch (error) {
    console.error('Cannot get the recipe:', error);
  }

};

const submitForm = async () => {
  const response = await fetch('/api/recipe/' + recipeId.value, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: recipeName.value, content: recipe.value }),
  });

  console.log(recipeId.value)

  if (response.ok) {
    alert('Uploaded Recipe');
  } else if (response.status >= 400 && response.status < 500) {
    alert('Client Error: ${response.statusText}');
  } else if (response.status >= 500) {
    alert('Server Error: Please try again later.');
  } else {
    alert('Unexpected Error: Can\'t Upload Recipe');
  }
};

// Fetch when component is mounted
onMounted(getRecipe);

// Watch for route changes and refetch
watch(() => route.params.id, (newId) => {
  recipeId.value = newId;
  getRecipe();
});

const input = computed(() => `# ${recipeName.value} \n ${recipe.value}`);

const output = computed(() => marked(input.value))

</script>

<template>
  <h1 class="title text-center pt-10 pb-1 mx-100 ">Edit a Recipe</h1>

  <div class="flex flex-row items-start justify-center px-30 h-screen">

    <div class="card-container text-center w-dvh">

      <div class="card-content">

        <form @submit.prevent="submitForm" class="flex flex-col gap-4">
          <p class="flex flex-col">
            <label for="name" class="f-label">Name</label>
            <br>
            <input v-model="recipeName" type="text" name="name" class="f-input" required />
          </p>

          <p class="flex flex-col">
            <label for="content" class="f-label">Write here</label>
            <br>
            <textarea v-model="recipe" name="content" rows="10"
              class="border-2 border-dashed border-amber-700 bg-white p-5 text-black font-roboto-mono"
              required></textarea>
          </p>

          <div class="relative inline-block mt-10">
            <button type="submit" class="taped-button">Submit</button>
            <div class="tape tape-top"></div>
            <div class="tape tape-bottom"></div>
          </div>

        </form>

      </div>

    </div>

    <!-- <div class="bg-amber-100 text-amber-900 text-left ml-14 px-5 py-3 w-dvh h-6/7 overflow-auto overflow-x-auto border border-gray-500"> -->
    <!--   <div class="output prose lg:prose-xl" v-html="output"></div> -->
    <!-- </div> -->

    <div class="card-container mx-5">
      <div class="card-content">

        <div class="output prose min-h-100 max-h-150 min-w-100 overflow-auto w-dvh 
                 lg:prose-li:marker:text-amber-900 lg:prose-headings:text-amber-900 
                 lg:prose-headings:font-[Monomakh] lg:text-amber-900" v-html="output">
        </div>

      </div>
    </div>

  </div>

</template>
