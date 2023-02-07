#version 330 core
out vec4 FragColor;


uniform vec3 cubeColor;

void main()
{
  // FragColor = vec4(1.0);


  // ATTEMPT: trying to change the light cube color as we change the light color
  FragColor = vec4(cubeColor, 1.0);
  
}
