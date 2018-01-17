<template>
	<div class="form-group document-language">

		<label for="document-language" class="form-control-label">Language</label>
		<select
			type="text"
			name="language"
			class="form-control"
			v-model="document.language"
		>
			<option v-for="(language, i) in languages" :key="i" :value="language.code">
				{{ language.flag }} {{ language.name }}
			</option>

		</select>

	</div>
</template>

<script lang="babel">

	import Accessors from '../../Mixins/accessors';

	export default {
		name: "LanguageField",
		mixins: [Accessors],
		computed: {
			languages() {
				return this.$store.state.languages;
			}
		}
	};

</script>
