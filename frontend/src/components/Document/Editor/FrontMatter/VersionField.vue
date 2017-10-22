<template>
	<div class="form-group document-version">

		<label for="document-version" class="form-control-label">Version</label>
		<input
			id="tag-editor"
			type="text"
			name="version"
			class="form-control"
			v-model="document.version"
		/>

	</div>
</template>

<script lang="babel">

	import Accessors from '../../../Mixins/accessors';


	export default {
		name: "VersionField",
		mixins: [Accessors]
	}
</script>
