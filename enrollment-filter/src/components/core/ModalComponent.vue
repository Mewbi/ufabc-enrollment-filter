<script setup lang="ts">
import { ref } from "vue";
import api from "../../services/api";
import LoadingIcon from "../icons/LoadingIcon.vue";
import { closeModal } from "jenesius-vue-modal";
import { getDocument } from "pdfjs-dist";
import * as pdfjs from "pdfjs-dist";

// @ts-ignore: import error
import PdfjsWorker from "pdfjs-dist/build/pdf.worker";
import { TextItem } from "pdfjs-dist/types/src/display/api";

pdfjs.GlobalWorkerOptions.workerSrc = PdfjsWorker;

const name = ref("");
const url = ref("");
const pdfFile = ref<File>();

const props = defineProps({
  apiStatus: {
    type: Boolean,
    required: true,
  },
  fileOptions: {
    type: Map,
    required: true,
  },
  localFiles: {
    type: Map,
    required: true,
  },
});

const newEnrollment = async () => {
  isDisabledModal.value = true;

  if (props.apiStatus) {
    getFromApi();
    return;
  }

  await getFromPDF();
  isDisabledModal.value = false;
  closeModal();
};

const getFromPDF = async () => {
  if (pdfFile.value === undefined) {
    alert("Faça o upload de um arquivo PDF válido.");
    isDisabledModal.value = false;

    return;
  }

  if (name.value === undefined || name.value.length === 0) {
    alert("Insira um nome válido.");
    isDisabledModal.value = false;
    return;
  }

  const buffer = await pdfFile.value.arrayBuffer();
  const load = getDocument({ data: buffer });
  const regex = /([0-9]+)\s([A-Z0-9]+-[0-9]{2}(SA|SB))\s(.*)/;
  const csvRows: string[][] = []; // I need use array because I cannot pass string as reference

  props.localFiles.set(name.value, csvRows);
  props.fileOptions.set(name.value, name.value);

  load.promise.then(
    (pdf) => {
      for (let i = 1; i <= pdf.numPages; i++) {
        pdf.getPage(i).then((page) => {
          const textLoad = page.getTextContent();
          textLoad.then(
            (txt) => {
              let line = "";
              let posY = 0;
              txt.items.forEach((content) => {
                const item = content as TextItem;
                const actualPosY = item.transform[5];

                if (posY != actualPosY) {
                  matchLineRegex(line);
                  line = "";
                  posY = actualPosY;
                }
                line += item.str;
              });
            },
            (err) => {
              console.log(err);
            }
          );
        });
      }
    },
    (err) => {
      console.log(err);
    }
  );

  const matchLineRegex = (line: string) => {
    const elements = line.match(regex);
    if (elements && elements.length >= 5) {
      csvRows.push([elements[1].toString(), elements[2].toString(), `"${elements[4].toString()}"`]);
    }
  };
};

const onFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    pdfFile.value = target.files[0];
  }
};

const getFromApi = async () => {
  const response = await api.newEnrollment(name.value, url.value);

  if (api.isErrorResponse(response)) {
    if ("error" in response) {
      if ("message" in response && response.message) alert(response.message);
    }
    isDisabledModal.value = false;
    return;
  }

  props.fileOptions.set(response.name, response.id);
  isDisabledModal.value = false;
  closeModal();
};

const isDisabledModal = ref<boolean>(false);
</script>

<template>
  <div
    class="flex rounded modal-window bg-white h-2/4 min-h-[450px] w-full max-w-xl justify-center items-center"
  >
    <div
      class="flex flex-col justify-between space-y-5 items-center w-full"
      v-bind:class="{ 'opacity-40': isDisabledModal }"
    >
      <div class="space-y-2 pl-5 pr-5">
        <p class="font-bold text-lg">Upload novo Arquivo</p>
      </div>
      <div
        class="space-y-2 pl-5 pr-5"
        v-bind:class="{ 'space-y-5': !props.apiStatus }"
      >
        <p v-if="props.apiStatus">
          Faça o upload de um novo arquivo PDF informando o nome do arquivo e
          sua respectiva URL.
        </p>
        <p v-else>
          Faça o upload de um novo arquivo PDF informando e informe o nome do
          arquivo.
        </p>
        <p>
          Clique
          <a
            href="https://prograd.ufabc.edu.br/pdf/ajuste_2024_2_matriculas_deferidas.pdf"
            target="_blank"
            class="font-bold"
            >aqui</a
          >
          para visualizar o formato do arquivo esperado.
        </p>
        <p v-if="props.apiStatus">
          Uma vez feito o upload não será possível alterar o nome informado.
          Escolha um bom nome =)
        </p>
        <p v-else>
          Caso a página seja recarregada será necessário refazer o upload do
          arquivo.
        </p>
        <p v-if="props.apiStatus">
          São aceitos apenas links do subdomínio <b>prograd.ufabc.edu.br</b>.
        </p>
      </div>

      <div class="w-full space-y-2">
        <div class="flex w-full pl-5 pr-5">
          <label
            class="border pl-4 w-20 h-full rounded-s-md p-2.5 bg-gray-200 font-bold text-sm"
            ><b>Nome</b></label
          >
          <input
            :disabled="isDisabledModal"
            type="text"
            v-model="name"
            placeholder="Matrículas Deferidas Pós Reajuste 2024.2"
            class="border ps-3 text-sm rounded-e-md w-full p-2.5 t-bold"
          />
        </div>

        <div v-if="props.apiStatus" class="flex w-full pl-5 pr-5">
          <label
            class="border pl-4 w-20 h-full rounded-s-md p-2.5 bg-gray-200 font-bold text-sm"
            ><b>URL</b></label
          >
          <input
            :disabled="isDisabledModal"
            type="text"
            v-model="url"
            placeholder="https://prograd.ufabc.edu.br/pdf/ajuste_2024_2_matriculas_deferidas.pdf"
            class="border ps-3 text-sm rounded-e-md w-full p-2.5 t-bold"
          />
        </div>
        <div class="flex w-full pl-5 pr-5">
          <label
            class="border pl-4 w-20 h-full rounded-s-md p-2.5 bg-gray-200 font-bold text-sm"
            ><b>Arquivo</b></label
          >
          <input
            :disabled="isDisabledModal"
            type="file"
            @change="onFileChange"
            class="border ps-3 text-sm rounded-e-md w-full p-1.5 t-bold"
          />
        </div>
      </div>
      <div>
        <button
          :disabled="isDisabledModal"
          class="border bg-green-800 text-white font-bold inline-flex items-center justify-center whitespace-nowrap rounded text-sm h-10 px-4 py-2"
          @click="newEnrollment"
        >
          Upload
        </button>
      </div>
      <div
        v-if="isDisabledModal"
        role="status"
        class="absolute -translate-x-1/2 -translate-y-1/2 top-2/4 left-1/2"
      >
        <LoadingIcon />
      </div>
    </div>
  </div>
</template>

<style></style>
