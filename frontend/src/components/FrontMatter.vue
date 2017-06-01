<template>

	<div class="frontmatter-fields">
		<div class="form-group">
			<label for="title">Title</label>
			<input name="title" class="form-control" v-model="document.title"/>
		</div>

		<div class="form-group">
			<label for="synopsis">Synopsis</label>
			<textarea name="synopsis" class="form-control" v-model="document.synopsis"/>
		</div>

		<div class="form-group">
			<label for="author">Author</label>
			<input name="author" class="form-control" v-model="document.author"/>
		</div>

		<div class="form-group">
			<label for="tags">Tags</label>
			<input name="tags" class="form-control" v-model="document.tags"/>
		</div>
	</div>

</template>

<script lang="babel">
	export default {
		name: "FrontMatter",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		}
	}
</script>