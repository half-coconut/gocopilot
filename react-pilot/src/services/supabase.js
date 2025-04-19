import { createClient } from "@supabase/supabase-js";
export const supabaseUrl = "https://nvbtdjgdbhgsgccuwkap.supabase.co";
// const supabaseKey = process.env.SUPABASE_KEY;
const supabaseKey =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im52YnRkamdkYmhnc2djY3V3a2FwIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA4MDkyNTIsImV4cCI6MjA1NjM4NTI1Mn0.Eq7kd3c9YjJBgKzTJPbVSVc40GEk2qbRzY3cYnry9kg";
export const supabase = createClient(supabaseUrl, supabaseKey);

export default supabase;
