<template>
	<div class="form-group">
		<label for="title">Title</label>
		<input name="title" class="form-control" v-model="document.title"/>
	</div>
</template>

<script lang="babel">
	export default {
		name: "TitleField",
		computed: {
			document() {
				console.log("ad", this.$store.state.activeDocument)
				return this.$store.state.activeDocument;
			}
		},
	};
</script>