#version 120

varying vec3 normal, lightDir;
varying vec3 tcoord;
uniform sampler2D textureStone, textureGrass;

void main() {
	vec4 texelStone, texelGrass;
	vec4 texel;
	vec3 n, l;
	float ndotl;

	n = normalize(normal);
	l = normalize(lightDir);
	ndotl = max(dot(n, l), 0.0);

	texelStone = texture2D(textureStone, tcoord.xy*10);
	texelGrass = texture2D(textureGrass, tcoord.xy*10);

	//gl_FragColor = vec4(tcoord.x, tcoord.y, tcoord.z*10 , 1.0) ;
	gl_FragColor = mix(texelGrass, texelStone, tcoord.z*5)*(ndotl*0.7+0.3);
}

