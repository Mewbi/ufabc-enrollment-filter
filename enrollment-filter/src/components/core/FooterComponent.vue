<script setup lang="ts">
import GithubIcon from "../icons/GithubIcon.vue";
</script>

<template>
  <footer class="w-full max-w-md">
    <div class="flex justify-center mt-8 space-x-2">
      <a
        target="_blank"
        href="https://github.com/Mewbi/ufabc-enrollment-filter"
        class="hover:underline hover:text-gray-900 flex gap-2"
      >
        <GithubIcon />
      </a>
      <p class="text-gray-600">- Feito por</p>
      <a
        target="_blank"
        href="https://github.com/Mewbi/"
        class="hover:underline hover:text-gray-900 text-gray-600 flex gap-2"
      >
        Felipe Fernandes
      </a>
    </div>

    <div>
      <hr class="my-2 border-gray-400 sm:mx-auto lg:my-2" />
      <span class="block text-sm text-gray-500 sm:text-center">
        © {{ new Date().getFullYear() }}
        <RouterLink :to="{ name: 'home' }" class="hover:underline">
          - UFABC Enrollment Filter</RouterLink
        >
      </span>
    </div>
  </footer>
</template>

<style scoped></style>
