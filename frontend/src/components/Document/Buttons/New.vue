<template>
	<router-link :to="{name: 'document_new', params: {directory: directoryPath}}" class="btn btn-sm btn-secondary">
		New document
	</router-link>
</template>

<script lang="babel">

	import Accessors from '../../Mixins/accessors';

	export default {
		name: "NewDocumentButton",
		mixins: [Accessors],
		props: ["directoryPath"]
	};

</script>