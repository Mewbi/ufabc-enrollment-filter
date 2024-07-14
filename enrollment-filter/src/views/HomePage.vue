<script setup lang="ts">
import Papa from "papaparse";
import Multiselect from "vue-multiselect";
import { saveAs } from "file-saver";
import "vue-good-table-next/dist/vue-good-table-next.css";
import { VueGoodTable } from "vue-good-table-next";
import { ref, onMounted, computed, watch } from "vue";
import LoadingIcon from "../components/icons/LoadingIcon.vue";
import FooterComponent from "../components/core/FooterComponent.vue";
import ModalComponent from "../components/core/ModalComponent.vue";
import { openModal, container } from "jenesius-vue-modal";
import api from "../services/api";

interface Column {
  label: string;
  field: string;
  type: string;
}

interface Row {
  ra: string;
  code: string;
  name: string;
}

interface Class {
  code: string;
  name: string;
}

const rows = ref<Row[]>([]);
const columns = ref<Column[]>([
  {
    label: "RA",
    field: "ra",
    type: "string",
  },
  {
    label: "Código Turma",
    field: "code",
    type: "string",
  },
  {
    label: "Nome da Disciplina",
    field: "name",
    type: "string",
  },
]);

const fileOptions = ref<Map<string, string>>();

const file = ref<string>();
const apiStatus = ref<boolean>(false);

const getApiStatus = async () => {
  apiStatus.value = await api.getApiStatus();

  if (apiStatus.value) {
    const res = await api.getEnrollmentList();

    if (api.isErrorResponse(res)) {
      if ("error" in res) {
        if ("message" in res && res.message) console.log(res.message);
      }
      return;
    }

    fileOptions.value = new Map<string, string>(
      res.map((enrollment) => [enrollment.name, enrollment.id])
    );

    // Set the first item of options as default value
    file.value = fileOptions.value.entries().next().value[1];
    return;
  }

  fileOptions.value = new Map<string, string>([
    [
      "Matrículas Deferidas Pós Reajuste 2024.2",
      "reajuste_2024_2_matriculas_deferidas.csv",
    ],
    [
      "Matrículas Deferidas Pós Ajuste 2024.2",
      "ajuste_2024_2_matriculas_deferidas.csv",
    ],
  ]);

  // Set the first item of options as default value
  file.value = fileOptions.value.entries().next().value[1];
};

const onFileLoaded = (results: Papa.ParseResult<Row>) => {
  rows.value = results.data;
  isDisabled.value = false;
};

const fetchCSV = async () => {
  isDisabled.value = true;

  if (apiStatus.value) {

    const response = await api.getEnrollment(`${file.value}`) 
    if (api.isErrorResponse(response)) {
      if ("error" in response) {
        if ("message" in response && response.message) 
          console.log(response.message);
      }
      return;
    }

    parseCSV(response);
    return;
  }

  try {
    const baseURL = `/ufabc-enrollment-filter/`;
    const response = await fetch(`${baseURL}${file.value}`);
    const csv = await response.text();
    parseCSV(csv);
  } catch (error) {
    console.error("Failed to fetch CSV:", error);
  }
};

const parseCSV = (csv: string) => {
  Papa.parse(csv, {
      header: true,
      transformHeader: (_: string, index: number) => {
        return columns.value[index].field;
      },
      complete: onFileLoaded,
    });
}

onMounted(async () => {
  await getApiStatus();
  fetchCSV();
});

const searchValue = ref<Class>();
const searchValueCode = computed(() => {
  return searchValue.value?.code;
});

const searchOptions = computed(() => {
  const optionsMap = rows.value.reduce<Map<string, Class>>((map, row) => {
    if (!map.has(row.code) && row.code !== undefined)
      map.set(row.code, { code: row.code, name: row.name } satisfies Class);
    return map;
  }, new Map<string, Class>([["", { code: "", name: "Limpar busca" } satisfies Class]]));

  const options = Array.from(optionsMap.values());
  return options;
});

const searchLabel = ({ code, name }: Class) => {
  if (code !== undefined && code.length === 0) return `${name}`; // Show 'Limpar' label
  return `${code} - ${name}`;
};

const raFiltered = ref<string[][]>([]);
let auxFiltered = false;

watch(rows, () => {
  // If new csv was loaded, fill RAs
  raFiltered.value = [];
  rows.value.forEach((r) => raFiltered.value.push([r.ra]));
});

const search = (row: Row, _: Column, __: string, searchTerm: string) => {
  if (auxFiltered) {
    raFiltered.value = [];
    auxFiltered = false;
  }

  if (row.code === searchTerm || searchTerm.length === 0) {
    raFiltered.value?.push([row.ra]);
    return true;
  }
  return false;
};

const onSearch = (params: { searchTerm: string; rowCount: number }) => {
  if (params.searchTerm.length === 0) {
    searchValue.value = undefined; // Clear input
  }
  auxFiltered = true;
};

const downloadRA = () => {
  const csv = Papa.unparse({
    fields: ["RA"],
    data: raFiltered.value,
  });

  const blob = new Blob([csv], { type: "text/csv" });
  saveAs(blob, file.value);
};

const isDisabled = ref<boolean>(false);

const showModal = () => {
  openModal(ModalComponent, {
    apiStatus: apiStatus,
    fileOptions: fileOptions,
  });
};
</script>

<template>
  <div class="flex flex-col w-full min-h-screen justify-between">
    <div>
      <header
        class="bg-background border-b p-4 md:p-6 flex flex-col gap-4 sm:flex-row sm:items-center"
      >
        <div class="space-y-2">
          <div class="flex place-items-center">
            <img src="../assets/ufabc.png" class="w-12" />
            <h1 class="text-2xl font-bold ps-2">
              UFABC - Parser de Matrículas
            </h1>
          </div>
          <p class="text-muted-foreground">
            Filtre as planilhas de matrículas deferidas facilmente.
          </p>
        </div>

        <div class="flex items-center ml-auto">
          <label
            class="ps-4 p-2.5 h-10 rounded-s-md font-bold text-sm bg-gray-200"
          >
            Arquivo
          </label>
          <select
            v-model="file"
            class="h-10 w-full text-sm p-2.5 border bg-white"
          >
            <option
              v-for="[key, value] of fileOptions"
              :key="key"
              :value="value"
            >
              {{ key }}
            </option>
          </select>
          <button
            @click="fetchCSV"
            :disabled="isDisabled"
            v-bind:class="{ 'opacity-40': isDisabled }"
            class="bg-green-600 text-white font-bold inline-flex items-center justify-center whitespace-nowrap text-sm h-10 px-4 py-2"
          >
            Carregar
          </button>
          <button
            @click="showModal"
            :disabled="isDisabled"
            v-bind:class="{ 'opacity-40': isDisabled }"
            class="bg-green-800 text-white font-bold inline-flex items-center justify-center whitespace-nowrap rounded-e-md text-sm h-10 px-4 py-2"
          >
            Novo Arquivo
          </button>
        </div>
      </header>

      <main class="flex p-4 md:p-6">
        <container class="z-50" />
        <div class="rounded-md w-full">
          <div class="flex items-center ml-auto mb-3">
            <label
              class="ps-4 p-2.5 h-10 rounded-s-md font-bold text-sm bg-gray-200"
            >
              Pesquisa
            </label>
            <multiselect
              v-model="searchValue"
              class="multiselect text-clip bg-inherit w-full text-sm border h-10 p-2.5 z-30 bg-white text-gray-600"
              :options="searchOptions"
              :custom-label="searchLabel"
              placeholder="Pesquise por uma turma"
            ></multiselect>
            <button
              @click="downloadRA"
              :disabled="isDisabled"
              v-bind:class="{ 'opacity-40': isDisabled }"
              class="border bg-green-800 text-white font-bold inline-flex items-center justify-center whitespace-nowrap rounded-e-md text-sm h-10 px-4 py-2"
            >
              Baixar RAs
            </button>
          </div>

          <div class="relative w-full overflow-auto">
            <vue-good-table
              v-bind:class="{ 'opacity-40': isDisabled }"
              :columns="columns"
              :rows="rows"
              :pagination-options="{
                enabled: true,
                rowsPerPageLabel: 'Itens por página',
                nextLabel: 'Próximo',
                prevLabel: 'Anterior',
              }"
              :sort-options="{
                enabled: true,
              }"
              :search-options="{
                enabled: false,
                externalQuery: searchValueCode,
                searchFn: search,
              }"
              v-on:search="onSearch"
            >
            </vue-good-table>
            <div
              v-if="isDisabled"
              role="status"
              class="absolute -translate-x-1/2 -translate-y-1/2 top-2/4 left-1/2"
            >
              <LoadingIcon />
            </div>
          </div>
        </div>
      </main>
    </div>
    <!-- Footer section -->
    <div class="flex flex-col items-center text-center pb-3">
      <FooterComponent />
    </div>
  </div>
</template>

<style>
.multiselect {
  cursor: pointer;
}
.multiselect__option--highlight {
  background-color: rgba(0, 0, 0, 0.1);
}

.multiselect__option {
  height: 10px;
  background-color: #ffffff;
  line-height: 2;
  padding: 5px 10px;
}
.multiselect__option--highlight {
  background-color: #dff0d8;
  color: #3c763d;
}

.multiselect__content {
  background-color: #ffffff;
}
.multiselect__single,
.multiselect__tag {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.multiselect__tags {
  display: flex;
  flex-wrap: nowrap;
}

.multiselect__tag {
  max-width: 100%;
}

.multiselect__tags-wrap {
  display: flex;
  flex-wrap: nowrap;
  overflow: hidden;
}
</style>
