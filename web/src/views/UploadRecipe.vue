<script setup>
  import { marked } from 'marked'
  import { debounce } from 'lodash-es'
  import { ref, computed } from 'vue'; 

  const name = ref('//Name//');
  const content = ref('# Please Write in md: [Guid](https://www.markdownguide.org/basic-syntax/)');

  const submitForm = async () => {
    const response = await fetch('/api/uploadrecipe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name.value, content: content.value}),
    });   

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

  const input = computed(() => `# ${name.value} \n ${content.value}`);

  const output = computed(() => marked(input.value))


</script>

<template>

  <h1 class="title text-center pt-10 pb-1 mx-100 ">Upload a Recipe</h1>

  <div class="flex flex-row items-start justify-center px-30 h-screen">

    <div class="card-container text-center w-dvh" >

      <div class="card-content">

        <form @submit.prevent="submitForm" class="flex flex-col gap-4">
          <p class="flex flex-col">
            <label for="name" class="f-label">Name</label>
            <br>
            <input v-model="name" type="text" name="name" class="f-input" required />
          </p>

          <p class="flex flex-col">
            <label for="content" class="f-label">Write here</label>
            <br>
            <textarea v-model="content" name="content" rows="10" class="border-2 border-dashed border-amber-700 bg-white p-5 text-black font-roboto-mono" required></textarea>
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

        <div 
          class="output prose min-h-100 max-h-150 min-w-100 overflow-auto w-dvh 
                 lg:prose-li:marker:text-amber-900 lg:prose-headings:text-amber-900 
                 lg:prose-headings:font-[Monomakh] lg:text-amber-900" 
          v-html="output">
        </div>

      </div>
    </div>

  </div>

</template>
