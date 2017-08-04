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

		<TagField/>

	</div>

</template>

<script lang="babel">
	import TagField from "../components/TagField";

	export default {
		name: "FrontMatter",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		},
		components: {
			TagField
		}
	}
</script>