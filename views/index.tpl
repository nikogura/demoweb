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
	    <table width="100%">
	        <tr>
	            <th> Path </th>
	            <th> Purpose </th>
	        </tr>
	        <tr>
	            <td align="center"> <a href="/"> / </a></td>
	            <td align="center"> Unauthenticated Landing Page </td>
	        </tr>
	        <tr>
	            <td align="center"> <a href="/loggedin"> /loggedin </a></td>
	            <td align="center"> Authenticated Landing Page </td>
	        </tr>
	        <tr>
	            <td align="center"> <a href="/login"> /login </a></td>
	            <td align="center"> Login Page </td>
	        </tr>
	        <tr>
	            <td align="center"> <a href="/logout"> /logout </a></td>
	            <td align="center"> Logout Page </td>
	        </tr>
	    </table>
		<p>
			<a href="{{.Url}}" target="_blank">
				<button>{{.Label}}</button>
			</a>
		</p>
		<h2> {{.Content}} </h2>
	</body>
</html>
