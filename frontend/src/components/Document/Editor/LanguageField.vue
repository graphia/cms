<template>
	<div class="form-group document-language">

		<label for="document-language" class="form-control-label">Language</label>
		<select
			type="text"
			name="language"
			class="form-control"
			v-model="document.language"
		>
			<option value="en">English 🇬🇧</option>
			<option value="de">German 🇩🇪</option>
			<option value="fr">French 🇫🇷</option>

		</select>

	</div>
</template>

<script lang="babel">

	/*
	 *  TODO set default language to that which is set in the
	 *  CMS's config file
	 */

	import Accessors from '../../Mixins/accessors';

	export default {
		name: "LanguageField",
		mixins: [Accessors]
	};

</script>
