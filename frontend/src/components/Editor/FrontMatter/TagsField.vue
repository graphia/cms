<template>

	<div class="form-group">

		<label for="tags">Tags</label>
		<input id="tag-editor" type="tags" name="tags" class="form-control" v-model="document.tags"/>

	</div>

</template>

<script lang="babel">

	import TagsInput from "tags-input";

	export default {
		name: "Tags",
		data() {
			return {
				taggle: null
			};
		},
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		},
		mounted() {
			console.log("Creating");



			TagsInput(document.querySelector('input[type="tags"]'));

		}
	};
</script>

<style lang="scss">
.tags-input {
	/*
	display: inline-block;
	padding: 0 2px;
	background: #FFF;
	border: 1px solid #CCC;
	width: 16em;
	border-radius: 2px;
	box-shadow: inset 0 1px 2px rgba(0,0,0,0.1);
	*/
	width: 100%;
	background-color: #fff;
	background-image: none;
	background-clip: padding-box;
	border: 1px solid rgba(0,0,0,.15);
	border-radius: .25rem;
	box-shadow: inset 0 1px 1px rgba(0,0,0,.075);
	transition: border-color ease-in-out .15s,box-shadow ease-in-out .15s;
	padding: 0.5rem 0.75rem;

	.tag {
		display: inline-block;
		padding: .25em .4em;
		font-size: 75%;
		font-weight: 700;
		line-height: 1;
		color: #fff;
		text-align: center;
		white-space: nowrap;
		vertical-align: baseline;
		background-color: #0275d8;
		padding-right: .6em;
		padding-left: .6em;
		margin-right: 0.6em;
		border-radius: 10rem;

		&.selected {
			background-color: #777;
			border-color: #777;
			color: #EEE;
		}

		&.dupe {
			transform: scale3d(1.2,1.2,1.2);
			background-color: #FCC;
			border-color: #700;
		}
	}

	input {
		display: inline-block !important;
		padding: 3px;
		margin: 0 !important;
		background: none !important;
		border: none !important;
		box-shadow: none !important;
		font: inherit !important;
		font-size: 100% !important;
		outline: none !important;
	}

	.selected ~ input {
		opacity: 0.3;
	}
}
</style>