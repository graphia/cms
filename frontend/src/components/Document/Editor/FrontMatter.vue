<template>

	<div class="frontmatter-fields">

		<TitleField/>

		<div class="form-group">
			<label for="synopsis">Synopsis</label>
			<textarea name="synopsis" class="form-control" v-model="document.synopsis"/>
		</div>

		<div class="form-group">
			<label for="author">Author</label>
			<input name="author" class="form-control" v-model="document.author"/>
		</div>

		<VersionField/>
		<TagsField/>

	</div>

</template>

<script lang="babel">
	import TitleField from "../Editor/FrontMatter/TitleField";
	import TagsField from "../Editor/FrontMatter/TagsField";
	import VersionField from "../Editor/FrontMatter/VersionField";
	import Accessors from '../../Mixins/accessors';


	export default {
		name: "FrontMatter",
		mixins: [Accessors],
		components: {
			TitleField,
			TagsField,
			VersionField
		}
	}
</script>