<template>
  <div>
    <NoVncViewer />
    <VmOpts />
    <DiskOpts />
  </div>
</template>

<script setup>
import NoVncViewer from '@/components/NoVncViewer.vue'
import VmOpts from '@/components/VmOpts.vue'
import DiskOpts from '@/components/DiskOpts.vue'
</script>

<style scoped></style>
