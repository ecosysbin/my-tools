# import gradio as gr

# def greet(name, intensity):
#     return "Hello, " + name + "!!!" * int(intensity)

# demo = gr.Interface(
#     fn=greet,
#     inputs=["text", "slider"],
#     outputs=["text"],
# )

# demo.launch(share=True)

# import gradio as gr

# demo = gr.Interface(
#     fn=lambda x:x, 
#     inputs=gr.Image(type="filepath"), 
#     outputs=gr.Image()
# )
    
# demo.launch()
# import numpy as np
# import gradio as gr

# def sepia(input_img):
#     sepia_filter = np.array([
#         [0.393, 0.769, 0.189],
#         [0.349, 0.686, 0.168],
#         [0.272, 0.534, 0.131]
#     ])
#     sepia_img = input_img.dot(sepia_filter.T)
#     sepia_img /= sepia_img.max()
#     return sepia_img

# demo = gr.Interface(sepia, gr.Image(), "image")
# if __name__ == "__main__":
#     demo.launch()
from langchain_community.llms import Ollama
import gradio as gr

def sepia(input_img):
    print(type(input_img))
    print(input_img)
    llm = Ollama(model="gemma2:2b")
    result = llm.invoke(input_img)
    return result

demo = gr.Interface(sepia, gr.Image(type="filepath"), "text")
if __name__ == "__main__":
    demo.launch()