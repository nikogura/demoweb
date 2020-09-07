<!DOCTYPE html>
<html lang="en">
	<head>
		<style>
			h1 {text-align: center;}
			h2 {text-align: center;}
			p {text-align: center;}
		</style>
		<title> {{.Title}} </title>
	</head>
	<body>
		<p>
			<a href="{{.Url}}" target="_blank">
				<button>{{.Label}}</button>
			</a>
		</p>
		<h2> {{.Content}} </h2>
	</body>
</html>
