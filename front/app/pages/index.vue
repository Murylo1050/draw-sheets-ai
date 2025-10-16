<template>
  <div
    class="bg-black w-full h-[100vh] flex flex-col items-center justify-center"
  >
    <canvas
      ref="canvasRef"
      class="border border-gray-400 rounded w-[80%] h-[80%]"
    />
    <div class="mt-2 flex gap-2">
      <button class="px-3 py-1 bg-red-500 text-white rounded" @click="clear">
        Limpar
      </button>
      <button class="px-3 py-1 bg-green-500 text-white rounded" @click="save">
        Salvar
      </button>
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
        console.log("resposta", JSON.parse(res));
        // aqui você pode enviar para o backend
      }
    },
  },
};
</script>
