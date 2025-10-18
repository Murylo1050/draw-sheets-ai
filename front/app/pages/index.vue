<template>
  <div class="bg-[#3E4449] w-full h-[100vh] flex flex-col items-center">
    <canvas
      ref="canvasRef"
      class="border border-gray-400 rounded-2xl w-[80%] h-[70%] mt-5"
    />
    <div class="mt-2 flex gap-2">
      <button class="px-3 py-1 bg-red-500 text-white rounded" @click="clear">
        Limpar
      </button>
      <button class="px-3 py-1 bg-green-500 text-white rounded" @click="save">
        Salvar
      </button>
    </div>

    <div
      class="w-[80%] h-[25%] bg-[#5D6368] rounded-2xl flex flex-col justify-end p-4"
    >
      <div
        v-for="result in data"
        class="flex flex-col justify-end text-white overflow-y-scroll h-full"
      >
        <p>{{ result }}</p>
      </div>
      <input
        id=""
        placeholder="Dale"
        type="text"
        name=""
        class="w-full h-[20%] rounded-2xl bg-[#9AA0A6] ring-0 px-3"
      />
    </div>
  </div>
</template>

<script>
import SignaturePad from "signature_pad";

export default {
  name: "SignaturePad",
  data() {
    return {
      signaturePad: null,
      data: [],
    };
  },
  mounted() {
    if (import.meta.client && this.$refs.canvasRef) {
      const canvas = this.$refs.canvasRef;

      const resizeCanvas = () => {
        const ratio = Math.max(window.devicePixelRatio || 1, 1);
        canvas.width = canvas.offsetWidth * ratio;
        canvas.height = canvas.offsetHeight * ratio;
        canvas.getContext("2d").scale(ratio, ratio);
        if (this.signaturePad) this.signaturePad.clear();
      };

      // Observa mudanças no tamanho do container
      this.resizeObserver = new ResizeObserver(resizeCanvas);
      this.resizeObserver.observe(canvas.parentElement);

      resizeCanvas();

      this.signaturePad = new SignaturePad(canvas, {
        backgroundColor: "rgb(255, 255, 255)",
        penColor: "rgb(0, 0, 0)",
      });
    }
  },
  beforeUnmount() {
    this.resizeObserver?.disconnect();
  },
  methods: {
    clear() {
      this.signaturePad?.clear();
    },
    async save() {
      if (this.signaturePad && !this.signaturePad.isEmpty()) {
        const dataURL = this.signaturePad.toDataURL("image/png");

        const res = await $fetch("http://localhost:8080/generateImage", {
          method: "POST",
          body: {
            image_data: dataURL.split(",")[1],
          },
        });
        this.data = JSON.parse(res);
        // aqui você pode enviar para o backend
      }
    },
  },
};
</script>
