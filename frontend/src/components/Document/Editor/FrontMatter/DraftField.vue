<template>
	<div class="document-draft form-group">
		<input
			name="draft"
			type="checkbox"
			v-model="document.draft"
		/>

		<label for="draft" class="form-control-label">
			Draft
		</label>
	</div>
</template>

<script lang="babel">
	import Accessors from '../../../Mixins/accessors';

	export default {
		name: "PublishedField",
		mixins: [Accessors]
	};
</script>