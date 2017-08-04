<template>

	<div class="frontmatter-fields">

		<div class="form-group">
			<label for="synopsis">Synopsis</label>
			<textarea name="synopsis" class="form-control" v-model="document.synopsis"/>
		</div>

		<div class="form-group">
			<label for="author">Author</label>
			<input name="author" class="form-control" v-model="document.author"/>
		</div>

		<TagsField/>

	</div>

</template>

<script lang="babel">
	import TagsField from "../Editor/FrontMatter/TagsField";

	export default {
		name: "FrontMatter",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		},
		components: {
			TagsField
		}
	}
</script>