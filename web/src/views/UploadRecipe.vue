<script setup>
import { marked } from 'marked'
import { ref, computed, onMounted, onUnmounted } from 'vue';

const name = ref('//Name//');
const content = ref('# Please Write in md: [Guide](https://www.markdownguide.org/basic-syntax/)');
const activeTab = ref('edit'); // For mobile view
const isMobile = ref(false);

// Check if screen is mobile size
const checkScreenSize = () => {
  isMobile.value = window.innerWidth < 768;
};

// Setup event listeners for responsive behavior
onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});

const setActiveTab = (tab) => {
  activeTab.value = tab;
};

const submitForm = async () => {
  const response = await fetch('/api/uploadrecipe', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: name.value, content: content.value }),
  });

  if (response.ok) {
    alert('Uploaded Recipe');
  } else if (response.status >= 400 && response.status < 500) {
    alert(`Client Error: ${response.statusText}`);
  } else if (response.status >= 500) {
    alert('Server Error: Please try again later.');
  } else {
    alert('Unexpected Error: Can\'t Upload Recipe');
  }
};

const input = computed(() => `# ${name.value} \n ${content.value}`);
const output = computed(() => marked(input.value));
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <h1 class="text-2xl md:text-3xl font-bold text-center pt-10 pb-6">Upload a Recipe</h1>

    <!-- Mobile tabs -->
    <div class="md:hidden flex border-b border-amber-200 mb-6">
      <button @click="setActiveTab('edit')" class="flex-1 py-2 text-center font-medium"
        :class="activeTab === 'edit' ? 'text-amber-700 border-b-2 border-amber-700' : 'text-gray-500'">
        Edit
      </button>
      <button @click="setActiveTab('preview')" class="flex-1 py-2 text-center font-medium"
        :class="activeTab === 'preview' ? 'text-amber-700 border-b-2 border-amber-700' : 'text-gray-500'">
        Preview
      </button>
    </div>

    <div class="flex flex-col md:flex-row md:items-start md:justify-center">
      <!-- Form section -->
      <div class="w-full md:w-1/2 mb-6 md:mb-0 text-center" :class="{ 'hidden': activeTab === 'preview' && isMobile }">
        <div class="card-container">
          <div class="card-content">
            <form @submit.prevent="submitForm" class="flex flex-col gap-4">
              <p class="flex flex-col">
                <label for="name" class="f-label">Name</label>
                <br>
                <input v-model="name" type="text" name="name"
                  class="border-2 border-amber-700 bg-white p-2 text-black w-full" required />
              </p>

              <p class="flex flex-col">
                <label for="content" class="f-label">Write here</label>
                <br>
                <textarea v-model="content" name="content" rows="10"
                  class="border-2 border-dashed border-amber-700 bg-white p-5 text-black font-mono w-full"
                  required></textarea>
              </p>

              <div class="relative inline-block mt-10">
                <button type="submit" class="bg-amber-600 hover:bg-amber-700 text-white px-6 py-2 font-medium">
                  Submit
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- Preview section -->
      <div class="w-full md:w-1/2 md:ml-6" :class="{ 'hidden': activeTab === 'edit' && isMobile }">
        <div class="bg-amber-100 text-amber-900 text-left p-5 border border-gray-500 h-full overflow-auto">
          <div class="output prose prose-amber max-w-none min-h-[200px]
                   lg:prose-li:marker:text-amber-900 lg:prose-headings:text-amber-900 
                   lg:prose-headings:font-[Monomakh] lg:text-amber-900" v-html="output"></div>
        </div>
      </div>
    </div>
  </div>
</template>
